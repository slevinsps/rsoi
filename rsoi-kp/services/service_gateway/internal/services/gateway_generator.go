package api

import (
	"net/http"
	"services/utils"

	uuid "github.com/satori/go.uuid"
)

// GeneratorStart
func (h *Handler) GeneratorStart(rw http.ResponseWriter, r *http.Request) {
	const place = "GeneratorStart"
	utils.PrintDebug(place)
	var (
		code int
		body []byte
		err  error
	)

	if code, body, err = h.SendRequest(r, h.conf.ServiceGenerator.URL, "/start", "POST"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GeneratorStart " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GeneratorStop
func (h *Handler) GeneratorStop(rw http.ResponseWriter, r *http.Request) {
	const place = "GeneratorStop"
	utils.PrintDebug(place)
	var (
		code int
		body []byte
		err  error
	)

	if code, body, err = h.SendRequest(r, h.conf.ServiceGenerator.URL, "/stop", "POST"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GeneratorStop " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetDataByEquipmentUUID
func (h *Handler) GetDataByEquipmentUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetDataByEquipmentUUID"
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

	if code, body, err = h.SendRequest(r, h.conf.ServiceGenerator.URL, "/equipment/"+equipmentUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in UpdateEquipment " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// DeleteDataByEquipmentUUID
func (h *Handler) DeleteDataByEquipmentUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "DeleteDataByEquipmentUUID"
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

	if code, body, err = h.SendRequest(r, h.conf.ServiceGenerator.URL, "/equipment/"+equipmentUUID.String(), "DELETE"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in UpdateEquipment " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// ClearAllData
func (h *Handler) ClearAllData(rw http.ResponseWriter, r *http.Request) {
	const place = "ClearAllData"
	utils.PrintDebug(place)
	var (
		code int
		body []byte
		err  error
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceGenerator.URL, "/clear", "DELETE"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in ClearAllData " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}
