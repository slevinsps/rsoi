package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"services/config"
	"services/internal/database"
	"services/internal/models"
	"services/utils"

	"github.com/go-redis/redis/v7"
	uuid "github.com/satori/go.uuid"
)

type Handler struct {
	conf        *config.Configuration
	db          *database.DataBase
	redisClient *redis.Client
}

func NewHandler(db *database.DataBase, redisClient *redis.Client) *Handler {
	confPath := "conf.json"
	var (
		conf *config.Configuration
		err  error
	)
	if conf, err = config.Init(confPath); err != nil {
		return nil
	}
	return &Handler{conf: conf, db: db, redisClient: redisClient}
}

// CreateMonitor
func (h *Handler) CreateMonitor(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateMonitor"
	utils.PrintDebug(place)
	var (
		code         int
		body         []byte
		body_request []byte
		err          error
		userUUID     uuid.UUID
		monitor      models.Monitor
	)
	if userUUID, err = h.CheckUser(rw, r); err != nil {
		return
	}

	body_request, err = ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body_request, &monitor); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	monitor.UserUUID = userUUID
	body, err = json.Marshal(monitor)

	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if code, body, err = h.SendRequest(r, h.conf.ServiceMonitor.URL, "/create", "POST"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in RequestCreateMonitor " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// AddEquipment
func (h *Handler) AddEquipment(rw http.ResponseWriter, r *http.Request) {
	const place = "AddEquipment"
	utils.PrintDebug(place)
	var (
		code          int
		body          []byte
		err           error
		userUUID      uuid.UUID
		monitorUUID   uuid.UUID
		equipmentUUID uuid.UUID
	)
	if userUUID, err = h.CheckUser(rw, r); err != nil {
		http.Error(rw, "No such user", http.StatusBadRequest)
		return
	}

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

	var body_new []byte
	r_check, _ := http.NewRequest("GET", h.conf.ServiceMonitor.URL, bytes.NewBuffer(body_new))
	if code, body, err = h.SendRequest(r_check, h.conf.ServiceMonitor.URL, "/"+monitorUUID.String()+"/user/"+userUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in RequestCreateMonitor " + err.Error())
		return
	}

	if code != http.StatusOK {
		rw.WriteHeader(code)
		sendJSON(rw, body)
	}

	r_check, _ = http.NewRequest("GET", h.conf.ServiceMonitor.URL, bytes.NewBuffer(body_new))
	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/"+equipmentUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetEquipmentByEquipmentUUID " + err.Error())
		return
	}

	if code != http.StatusOK {
		rw.WriteHeader(code)
		sendJSON(rw, body)
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceMonitor.URL, "/"+monitorUUID.String()+"/add/"+equipmentUUID.String(), "POST"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in RequestCreateMonitor " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// DelEquipment
func (h *Handler) DelEquipment(rw http.ResponseWriter, r *http.Request) {
	const place = "DelEquipment"
	utils.PrintDebug(place)
	var (
		code          int
		body          []byte
		err           error
		userUUID      uuid.UUID
		monitorUUID   uuid.UUID
		equipmentUUID uuid.UUID
	)
	if userUUID, err = h.CheckUser(rw, r); err != nil {
		return
	}

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

	var body_new []byte
	r_check, _ := http.NewRequest("POST", h.conf.ServiceMonitor.URL, bytes.NewBuffer(body_new))
	if code, body, err = h.SendRequest(r_check, h.conf.ServiceMonitor.URL, "/"+monitorUUID.String()+"/user/"+userUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in RequestCreateMonitor " + err.Error())
		return
	}

	if code != http.StatusOK {
		rw.WriteHeader(code)
		sendJSON(rw, body)
	}

	r_check, _ = http.NewRequest("POST", h.conf.ServiceMonitor.URL, bytes.NewBuffer(body_new))
	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/"+equipmentUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetEquipmentByEquipmentUUID " + err.Error())
		return
	}

	if code != http.StatusOK {
		rw.WriteHeader(code)
		sendJSON(rw, body)
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceMonitor.URL, "/"+monitorUUID.String()+"/del/"+equipmentUUID.String(), "DELETE"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in RequestCreateMonitor " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetAllMonitorsByUserUUID
func (h *Handler) GetAllMonitorsByUserUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetAllMonitorsByUserUUID"
	utils.PrintDebug(place)
	var (
		code     int
		body     []byte
		err      error
		userUUID uuid.UUID
	)
	if userUUID, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceMonitor.URL, "/user/"+userUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in RequestCreateMonitor " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	fmt.Println(code)
	rw.WriteHeader(http.StatusOK)
	sendJSON(rw, body)
	return
}

// GetMonitorByMonitorUUID
func (h *Handler) GetMonitorByMonitorUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetMonitorByMonitorUUID"
	utils.PrintDebug(place)
	var (
		code        int
		body        []byte
		err         error
		monitorUUID uuid.UUID
		userUUID    uuid.UUID
	)
	if userUUID, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if monitorUUID, err = getUUID(r, "monitorUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceMonitor.URL, "/"+monitorUUID.String()+"/user/"+userUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in RequestCreateMonitor " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// DelMonitorByMonitorUUID
func (h *Handler) DelMonitorByMonitorUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "DelMonitorByMonitorUUID"
	utils.PrintDebug(place)
	var (
		code        int
		body        []byte
		err         error
		monitorUUID uuid.UUID
		userUUID    uuid.UUID
	)
	if userUUID, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if monitorUUID, err = getUUID(r, "monitorUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceMonitor.URL, "/"+monitorUUID.String()+"/user/"+userUUID.String(), "DELETE"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in RequestCreateMonitor " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}
