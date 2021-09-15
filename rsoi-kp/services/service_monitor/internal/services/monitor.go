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
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type Handler struct {
	rep  repInterface.RepInterface
	conf *config.Configuration
}

func NewHandler(rep repInterface.RepInterface) *Handler {
	confPath := "conf.json"
	var (
		conf *config.Configuration
		err  error
	)
	if conf, err = config.Init(confPath); err != nil {
		return nil
	}
	return &Handler{rep: rep, conf: conf}
}

func ExtractToken(r *http.Request) string {
	fmt.Println("ExtractToken")
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	fmt.Println("ExtractToken")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	fmt.Println("VerifyToken")
	tokenString := ExtractToken(r)
	fmt.Println(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		utils.PrintDebug(err.Error())
		return nil, err
	}
	fmt.Println("VerifyToken")
	return token, nil
}

func TokenValid(r *http.Request) (token *jwt.Token, err error) {
	fmt.Println("TokenValid")
	token, err = VerifyToken(r)
	if err != nil {
		return
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return
	}
	fmt.Println("TokenValid")
	return
}

// TokenAuthMiddleware
func (h *Handler) TokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			err error
		)
		_, err = TokenValid(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			utils.PrintDebug("No valid token in TokenAuthMiddleware" + err.Error())
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CreateMonitor
func (h *Handler) CreateMonitor(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateMonitor"
	utils.PrintDebug(place)
	var (
		err     error
		monitor models.Monitor
	)
	monitorUID := uuid.NewV4()

	if monitor, err = getMonitorFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getMonitorFromBody " + err.Error())
		return
	}

	if len(monitor.Name) == 0 {
		sendMessage(rw, "Monitor name is empty", http.StatusBadRequest, nil)
		return
	}

	if monitor.UserUUID == uuid.Nil {
		sendMessage(rw, "User uuid is empty", http.StatusBadRequest, nil)
		return
	}

	monitor.MonitorUUID = monitorUID

	if err = h.rep.CreateMonitor(monitor); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in CreateMonitor")
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusCreated)
	resBytes, _ := json.Marshal(monitor)
	sendJSON(rw, resBytes)

	return
}

func (h *Handler) GetAllMonitorsByUserUUID(rw http.ResponseWriter, r *http.Request) {

	const place = "GetAllUserMonitors"

	var (
		err       error
		monitors  []models.Monitor
		user_uuid uuid.UUID
	)

	rw.Header().Set("Content-Type", "application/json")

	if user_uuid, err = getUUID(r, "userUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if monitors, err = h.rep.GetAllMonitorsByUserUUID(user_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.PrintDebug("Error in GetAllMonitors")
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(monitors) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(monitors)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) GetMonitorByMonitorUUIDuserUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetMonitorByUUID"

	var (
		err              error
		monitorUUID      uuid.UUID
		userUUID         uuid.UUID
		monitor          models.Monitor
		checkFindMonitor bool
	)

	if monitorUUID, err = getUUID(r, "monitorUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if userUUID, err = getUUID(r, "userUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if monitor, checkFindMonitor, err = h.rep.GetMonitorByMonitorUUIDuserUUID(monitorUUID, userUUID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if checkFindMonitor {
		rw.WriteHeader(http.StatusOK)
		resBytes, _ := json.Marshal(monitor)
		sendJSON(rw, resBytes)
	} else {
		sendMessage(rw, "Can't find monitor by uuid: "+monitorUUID.String(), http.StatusNotFound, nil)
	}

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) DelMonitorByMonitorUUIDuserUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "DelMonitorByMonitorUUIDuserUUID"

	var (
		err              error
		monitorUUID      uuid.UUID
		userUUID         uuid.UUID
		checkFindMonitor bool
	)

	if monitorUUID, err = getUUID(r, "monitorUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if userUUID, err = getUUID(r, "userUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if _, checkFindMonitor, err = h.rep.GetMonitorByMonitorUUIDuserUUID(monitorUUID, userUUID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindMonitor {
		sendMessage(rw, "Can't find monitor by uuid: "+monitorUUID.String(), http.StatusNotFound, nil)
		return
	}

	if err = h.rep.DelMonitorByMonitorUUID(monitorUUID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	sendMessage(rw, "Monitor was deleted", http.StatusOK, nil)
	printResult(err, http.StatusOK, place)
	return
}

// AddEquipment
func (h *Handler) AddEquipment(rw http.ResponseWriter, r *http.Request) {
	const place = "AddEquipment"
	utils.PrintDebug(place)
	var (
		err           error
		monitorUUID   uuid.UUID
		equipmentUUID uuid.UUID
	)

	if monitorUUID, err = getUUID(r, "monitorUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if equipmentUUID, err = getUUID(r, "equipmentUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if err = h.rep.AddEquipment(monitorUUID, equipmentUUID); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in CreateEquipment " + err.Error())
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusCreated)

	return
}

func (h *Handler) DelEquipment(rw http.ResponseWriter, r *http.Request) {
	const place = "DelEquipment"

	var (
		err           error
		equipmentUUID uuid.UUID
		monitorUUID   uuid.UUID
	)

	if monitorUUID, err = getUUID(r, "monitorUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if equipmentUUID, err = getUUID(r, "equipmentUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if err = h.rep.DelEquipment(monitorUUID, equipmentUUID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
	return
}

func CreateServiceToken(login string) (*models.TokenDetails, error) {
	var err error

	td := &models.TokenDetails{}

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["source"] = "service"
	atClaims["service login"] = login
	atClaims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
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
