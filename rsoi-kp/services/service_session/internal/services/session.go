package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"services/config"
	"services/internal/models"
	"services/internal/repInterface"
	"services/utils"
	"strconv"

	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
)

type Handler struct {
	rep         repInterface.RepInterface
	conf        *config.Configuration
	redisClient *redis.Client
}

func NewHandler(rep repInterface.RepInterface, redisClient *redis.Client) *Handler {
	confPath := "conf.json"
	var (
		conf *config.Configuration
		err  error
	)
	if conf, err = config.Init(confPath); err != nil {
		return nil
	}
	return &Handler{rep: rep, conf: conf, redisClient: redisClient}
}

// CreateUser
func (h *Handler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateUser"

	var (
		err           error
		user          models.User
		hash_password string
		checkFindUser bool
	)
	userUID := uuid.NewV4()

	if user, err = getUserFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUserFromBody " + err.Error())
		return
	}

	user.UserUUID = userUID
	if hash_password, err = HashPassword(user.Password); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in HashPassword " + err.Error())
		return
	}
	user.Password = hash_password
	user.IsAdmin = false

	if checkFindUser, err = h.rep.CreateUser(user); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.PrintDebug("Error in CreateUser")
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	fmt.Println("checkFindUser ", checkFindUser)

	if !checkFindUser {
		rw.WriteHeader(http.StatusCreated)
		resBytes, _ := json.Marshal(user)
		sendJSON(rw, resBytes)

	} else {
		sendMessage(rw, "User with this login exists", http.StatusConflict, nil)
	}

	return
}

func (h *Handler) SignIn(rw http.ResponseWriter, r *http.Request) {
	const place = "SignIn"

	var (
		err           error
		checkFindUser bool
		ts            *models.TokenDetails
		user          models.User
		userSignIn    models.User
	)

	if userSignIn, err = getUserFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUserFromBody " + err.Error())
		return
	}

	if user, checkFindUser, err = h.rep.GetUserByLogin(userSignIn.Login); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindUser {
		sendMessage(rw, "User with this login not exists", http.StatusNotFound, nil)
		return
	}

	if !CheckPasswordHash(userSignIn.Password, user.Password) {
		sendMessage(rw, "Invalid password", http.StatusNotFound, nil)
		return
	}

	ts, err = CreateToken(user.UserUUID, user.IsAdmin)
	if err != nil {
		sendMessage(rw, "Error in creating token", http.StatusUnprocessableEntity, nil)
		return
	}
	saveErr := h.CreateAuth(user.UserUUID, ts)
	if saveErr != nil {
		sendMessage(rw, "Error in saving token", http.StatusUnprocessableEntity, nil)
		return
	}

	rw.WriteHeader(http.StatusOK)
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
		"is_admin":      strconv.FormatBool(user.IsAdmin),
	}
	resBytes, _ := json.Marshal(tokens)
	sendJSON(rw, resBytes)
}

func (h *Handler) Logout(rw http.ResponseWriter, r *http.Request) {
	au, err := ExtractTokenMetadata(r, "header")
	if err != nil {
		sendMessage(rw, "unauthorized", http.StatusUnauthorized, nil)
		return
	}
	deleted, delErr := h.DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 {
		sendMessage(rw, "unauthorized", http.StatusUnauthorized, nil)
		return
	}
	sendMessage(rw, "Successfully logged out", http.StatusOK, nil)
}

func (h *Handler) GetAllUsers(rw http.ResponseWriter, r *http.Request) {

	const place = "GetAllUsers"

	var (
		err    error
		users  []models.User
		userId uuid.UUID
	)

	tokenAuth, err := ExtractTokenMetadata(r, "header")
	if err != nil {
		sendMessage(rw, "unauthorized", http.StatusUnauthorized, nil)
		return
	}
	userId, err = h.FetchAuth(tokenAuth)
	if err != nil {
		sendMessage(rw, "unauthorized", http.StatusUnauthorized, nil)
		return
	}

	fmt.Println(userId)

	rw.Header().Set("Content-Type", "application/json")

	if users, err = h.rep.GetAllUsers(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.PrintDebug("Error in GetAllUsers")
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(users) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(users)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) GetUserByUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetUserByUUID"

	var (
		err           error
		userUID       uuid.UUID
		user          models.User
		checkFindUser bool
	)

	if userUID, err = getUserUID(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if user, checkFindUser, err = h.rep.GetUserByUUID(userUID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if checkFindUser {
		rw.WriteHeader(http.StatusOK)
		resBytes, _ := json.Marshal(user)
		sendJSON(rw, resBytes)
	} else {
		sendMessage(rw, "Can't find user by uuid: "+userUID.String(), http.StatusNotFound, nil)
	}

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) GetUserByToken(rw http.ResponseWriter, r *http.Request) {
	const place = "GetUserByToken"

	var (
		err     error
		userUID uuid.UUID
		ad      *models.AccessDetails
	)

	if ad, err = ExtractTokenMetadata(r, "body"); err != nil {
		sendMessage(rw, "unauthorized", http.StatusUnauthorized, nil)
		utils.PrintDebug("Error in ExtractTokenMetadata(r, body)")
		return
	}
	fmt.Println(ad.UserId)

	if userUID, err = h.FetchAuth(ad); err != nil {
		sendMessage(rw, "unauthorized", http.StatusUnauthorized, nil)
		utils.PrintDebug("Error in h.FetchAuth(ad)")
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	ad.UserId = userUID
	rw.WriteHeader(http.StatusOK)
	resBytes, _ := json.Marshal(ad)
	sendJSON(rw, resBytes)

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) ServiceRegister(rw http.ResponseWriter, r *http.Request) {
	const place = "ServiceRegister"

	var (
		err     error
		service models.Service
		td      *models.TokenDetails
	)

	if service, err = getServiceFromBody(r); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		utils.PrintDebug("Error in getServiceFromBody(r) " + err.Error())
		return
	}
	if service.Login != os.Getenv("GATEWAY_ID") || service.Password != os.Getenv("GATEWAY_SECRET") {
		http.Error(rw, "you do not have permission to access this service", http.StatusUnauthorized)
		utils.PrintDebug("Error in checksecrets")
		return
	}

	td, err = CreateServiceToken(service.Login)
	if err != nil {
		utils.PrintDebug("Error in creating token " + err.Error())
		sendMessage(rw, "Error in creating token", http.StatusUnprocessableEntity, nil)
		return
	}
	rw.WriteHeader(http.StatusOK)
	resBytes, _ := json.Marshal(td)
	sendJSON(rw, resBytes)
	return
}
