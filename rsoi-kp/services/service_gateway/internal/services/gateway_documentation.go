package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"services/utils"

	uuid "github.com/satori/go.uuid"
)

// CreateFile
func (h *Handler) CreateFile(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateFile"
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

	fmt.Println("1 ", equipmentModelUUID.String())

	var body_new []byte
	r_check, _ := http.NewRequest("POST", h.conf.ServiceEquipment.URL, bytes.NewBuffer(body_new))
	if code, body, err = h.SendRequest(r_check, h.conf.ServiceEquipment.URL, "/model/"+equipmentModelUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetEquipmentByEquipmentUUID " + err.Error())
		return
	}

	fmt.Println("2 ", code)

	if code != http.StatusOK {
		rw.WriteHeader(code)
		sendJSON(rw, body)
		return
	}

	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("3 ")

	req, err := http.NewRequest("POST", h.conf.ServiceDocumentation.URL+"/upload/"+equipmentModelUUID.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	if code, body, err = h.SendRequest(req, h.conf.ServiceDocumentation.URL, "/upload/"+equipmentModelUUID.String(), "POST"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetEquipmentByEquipmentUUID " + err.Error())
		return
	}
	// if err != nil {
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	fmt.Println("4 ")
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// UpdateFile
func (h *Handler) UpdateFile(rw http.ResponseWriter, r *http.Request) {
	const place = "UpdateFile"
	utils.PrintDebug(place)
	var (
		code int
		body []byte
		err  error
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceDocumentation.URL, "/update", "PUT"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in UpdateEquipment " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetAllFiles
func (h *Handler) GetAllFiles(rw http.ResponseWriter, r *http.Request) {
	const place = "GetAllFiles"
	utils.PrintDebug(place)
	var (
		code int
		body []byte
		err  error
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceDocumentation.URL, "/list", "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetAllEquipments " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetAllFilesByEquipmentUUID
func (h *Handler) GetAllFilesByEquipmentUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetAllFilesByEquipmentUUID"
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

	if code, body, err = h.SendRequest(r, h.conf.ServiceDocumentation.URL, "/equipment_model/"+equipmentModelUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetAllFilesByEquipmentUUID " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// GetFileByFileUUID
func (h *Handler) GetFileByFileUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetFileByFileUUID"
	utils.PrintDebug(place)
	var (
		code     int
		body     []byte
		err      error
		fileUUID uuid.UUID
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if fileUUID, err = getUUID(r, "fileUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceDocumentation.URL, "/"+fileUUID.String(), "GET"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in GetFileByFileUUID" + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}

// DelFileByFileUUID
func (h *Handler) DelFileByFileUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "DelFileByFileUUID"
	utils.PrintDebug(place)
	var (
		code     int
		body     []byte
		err      error
		fileUUID uuid.UUID
	)
	if _, err = h.CheckUser(rw, r); err != nil {
		return
	}

	if fileUUID, err = getUUID(r, "fileUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if code, body, err = h.SendRequest(r, h.conf.ServiceEquipment.URL, "/"+fileUUID.String(), "DELETE"); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in DelFileByFileUUID " + err.Error())
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(code)
	sendJSON(rw, body)
	return
}
