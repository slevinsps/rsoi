package api

import (
	"net/http"
	"services/utils"

	uuid "github.com/satori/go.uuid"
)

// CreateEquipment
func (h *Handler) CreateEquipment(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateEquipment"
	utils.PrintDebug(place)
	var (
		code int
		body []byte
		err  error
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/create", "POST"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in CreateEquipment " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// UpdateEquipment
func (h *Handler) UpdateEquipment(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdateEquipment"
	utils.PrintDebug(place)
	var (
		code int
		body []byte
		err  error
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/update", "PUT"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in UpdateEquipment " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetAllEquipments
func (h *Handler) GetAllEquipments(rw http.ResponseWriter, r *http.Request) {
	const place = "GetAllEquipments"
	utils.PrintDebug(place)
	var (
		code int
		body []byte
		err  error
	)

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/list", "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetAllEquipments " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetEquipmentsByMonitorUUID
func (h *Handler) GetEquipmentsByMonitorUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetEquipmentsByMonitorUUID"
	utils.PrintDebug(place)
	var (
		code        int
		body        []byte
		err         error
		monitorUUID uuid.UUID
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if monitorUUID, err = getUUID(r, "monitorUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/list/"+monitorUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetEquipmentsByMonitorUUID " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetNotAddedEquipmentsByMonitorUUID
func (h *Handler) GetNotAddedEquipmentsByMonitorUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetNotAddedEquipmentsByMonitorUUID"
	utils.PrintDebug(place)
	var (
		code        int
		body        []byte
		err         error
		monitorUUID uuid.UUID
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if monitorUUID, err = getUUID(r, "monitorUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/list/"+monitorUUID.String()+"/notadded", "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetEquipmentsByMonitorUUID " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetEquipmentByEquipmentUUID
func (h *Handler) GetEquipmentByEquipmentUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetEquipmentByEquipmentUUID"
	utils.PrintDebug(place)
	var (
		code          int
		body          []byte
		err           error
		equipmentUUID uuid.UUID
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if equipmentUUID, err = getUUID(r, "equipmentUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/"+equipmentUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetEquipmentByEquipmentUUID " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// DelEquipmentByEquipmentUUID
func (h *Handler) DelEquipmentByEquipmentUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "DelEquipmentByEquipmentUUID"
	utils.PrintDebug(place)
	var (
		code          int
		body          []byte
		err           error
		equipmentUUID uuid.UUID
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if equipmentUUID, err = getUUID(r, "equipmentUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/"+equipmentUUID.String(), "DELETE"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in RequestCreateMonitor " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}
