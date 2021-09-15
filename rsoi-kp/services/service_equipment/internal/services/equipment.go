package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"services/config"
	"services/constants"
	"services/internal/models"
	"services/internal/repInterface"
	"services/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

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

// CreateEquipment
func (h *Handler) CreateEquipment(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateEquipment"
	utils.PrintDebug(place)
	var (
		err                     error
		equipment               models.Equipment
		checkFindEquipmentModel bool
	)
	equipmentUID := uuid.NewV4()

	if equipment, err = getEquipmentFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getEquipmentFromBody " + err.Error())
		return
	}

	if len(equipment.Name) == 0 {
		sendMessage(rw, "Equipment name is empty", http.StatusBadRequest, nil)
		return
	}

	if _, checkFindEquipmentModel, err = h.rep.GetEquipmentModelByEquipmentModelUUID(equipment.EquipmentModelUUID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindEquipmentModel {
		sendMessage(rw, "Can't find equipment model by uuid: "+equipment.EquipmentModelUUID.String(), http.StatusNotFound, nil)
	}

	switch equipment.Status {
	case
		constants.EQUIPMENT_STATUS_ACTIVE,
		constants.EQUIPMENT_STATUS_INACTIVE,
		constants.EQUIPMENT_STATUS_DELETED:
	default:
		equipment.Status = constants.EQUIPMENT_STATUS_ACTIVE
	}

	equipment.EquipmentUUID = equipmentUID

	if err = h.rep.CreateEquipment(equipment); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in CreateEquipment " + err.Error())
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusCreated)
	resBytes, _ := json.Marshal(equipment)
	sendJSON(rw, resBytes)

	return
}

// UpdateEquipment
func (h *Handler) UpdateEquipment(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdateEquipment"
	utils.PrintDebug(place)
	var (
		err                error
		equipment_new      models.Equipment
		equipment_old      models.Equipment
		equipment_res      models.Equipment
		checkFindEquipment bool
	)

	if equipment_new, err = getEquipmentFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getEquipmentFromBody " + err.Error())
		return
	}
	equipment_uuid := equipment_new.EquipmentUUID

	if equipment_old, checkFindEquipment, err = h.rep.GetEquipmentByEquipmentUUID(equipment_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if !checkFindEquipment {
		sendMessage(rw, "Can't find equipment by uuid: "+equipment_uuid.String(), http.StatusNotFound, nil)
	}

	if equipment_new.Name == "" {
		equipment_new.Name = equipment_old.Name
	}
	if equipment_new.EquipmentModelUUID == uuid.Nil {
		equipment_new.EquipmentModelUUID = equipment_old.EquipmentModelUUID
	}

	switch equipment_new.Status {
	case
		constants.EQUIPMENT_STATUS_ACTIVE,
		constants.EQUIPMENT_STATUS_INACTIVE,
		constants.EQUIPMENT_STATUS_DELETED:
	default:
		equipment_new.Status = equipment_old.Status
	}

	if equipment_res, err = h.rep.UpdateEquipment(equipment_new); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in CreateEquipment " + err.Error())
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusCreated)
	resBytes, _ := json.Marshal(equipment_res)
	sendJSON(rw, resBytes)

	return
}

func (h *Handler) GetAllEquipments(rw http.ResponseWriter, r *http.Request) {

	const place = "GetAllUserEquipments"

	var (
		err        error
		equipments []models.Equipment
	)

	rw.Header().Set("Content-Type", "application/json")

	if equipments, err = h.rep.GetAllEquipments(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.PrintDebug("Error in GetAllEquipments")
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(equipments) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(equipments)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) GetEquipmentByEquipmentUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetEquipmentByUUID"

	var (
		err                error
		equipment_uuid     uuid.UUID
		equipment          models.Equipment
		checkFindEquipment bool
	)

	if equipment_uuid, err = getUUID(r, "equipmentUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if equipment, checkFindEquipment, err = h.rep.GetEquipmentByEquipmentUUID(equipment_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if checkFindEquipment {
		rw.WriteHeader(http.StatusOK)
		resBytes, _ := json.Marshal(equipment)
		sendJSON(rw, resBytes)
	} else {
		sendMessage(rw, "Can't find equipment by uuid: "+equipment_uuid.String(), http.StatusNotFound, nil)
	}

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) DelEquipmentByEquipmentUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "DelEquipmentByEquipmentUUID"

	var (
		err                error
		equipment_uuid     uuid.UUID
		checkFindEquipment bool
	)

	if equipment_uuid, err = getUUID(r, "equipmentUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if _, checkFindEquipment, err = h.rep.GetEquipmentByEquipmentUUID(equipment_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindEquipment {
		sendMessage(rw, "Can't find equipment by uuid: "+equipment_uuid.String(), http.StatusNotFound, nil)
		return
	}

	if err = h.rep.DelEquipmentByEquipmentUUID(equipment_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	sendMessage(rw, "Equipment was deleted", http.StatusOK, nil)
	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) GetEquipmentsByMonitorUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetEquipmentByUUID"

	var (
		err         error
		monitorUUID uuid.UUID
		equipments  []models.Equipment
	)

	if monitorUUID, err = getUUID(r, "monitorUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if equipments, err = h.rep.GetEquipmentsByMonitorUUID(monitorUUID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(equipments) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(equipments)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) GetNotAddedEquipmentsByMonitorUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetNotAddedEquipmentsByMonitorUUID"

	var (
		err         error
		monitorUUID uuid.UUID
		equipments  []models.Equipment
	)

	if monitorUUID, err = getUUID(r, "monitorUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if equipments, err = h.rep.GetNotAddedEquipmentsByMonitorUUID(monitorUUID); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(equipments) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(equipments)
		sendJSON(rw, resBytes)
	}

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
