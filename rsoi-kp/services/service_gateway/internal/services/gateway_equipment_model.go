package api

import (
	"net/http"
	"services/utils"

	uuid "github.com/satori/go.uuid"
)

// CreateEquipmentModel
func (h *Handler) CreateEquipmentModel(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateEquipmentModel"
	utils.PrintDebug(place)
	var (
		code int
		body []byte
		err  error
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/model/create", "POST"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in CreateEquipmentModel " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetAllEquipmentModels
func (h *Handler) GetAllEquipmentModels(rw http.ResponseWriter, r *http.Request) {
	const place = "GetAllEquipmentModels"
	utils.PrintDebug(place)
	var (
		code int
		body []byte
		err  error
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/model/list", "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetAllEquipmentModels " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetEquipmentModelByEquipmentModelUUID
func (h *Handler) GetEquipmentModelByEquipmentModelUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetEquipmentModelByEquipmentModelUUID"
	utils.PrintDebug(place)
	var (
		code               int
		body               []byte
		err                error
		equipmentModelUUID uuid.UUID
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if equipmentModelUUID, err = getUUID(r, "equipmentModelUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/model/"+equipmentModelUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetEquipmentModelByEquipmentModelUUID " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// DelEquipmentModelByEquipmentModelUUID
func (h *Handler) DelEquipmentModelByEquipmentModelUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "DelEquipmentModelByEquipmentModelUUID"
	utils.PrintDebug(place)
	var (
		code               int
		body               []byte
		err                error
		equipmentModelUUID uuid.UUID
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if equipmentModelUUID, err = getUUID(r, "equipmentModelUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/model/"+equipmentModelUUID.String(), "DELETE"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in RequestCreateMonitor " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}
