package api

import (
	"encoding/json"
	"net/http"
	"services/internal/models"
	"services/utils"

	uuid "github.com/satori/go.uuid"
)

// CreateEquipmentModel
func (h *Handler) CreateEquipmentModel(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateEquipmentModel"
	utils.PrintDebug(place)
	var (
		err            error
		equipmentModel models.EquipmentModel
	)
	equipmentModelUID := uuid.NewV4()

	if equipmentModel, err = getEquipmentModelFromBody(r); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getEquipmentModelFromBody " + err.Error())
		return
	}

	if len(equipmentModel.Name) == 0 {
		sendMessage(rw, "EquipmentModel name is empty", http.StatusBadRequest, nil)
		return
	}

	equipmentModel.EquipmentModelUUID = equipmentModelUID

	if err = h.rep.CreateEquipmentModel(equipmentModel); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		utils.PrintDebug("Error in CreateEquipmentModel " + err.Error())
		return
	}
	rw.Header().Set("Content-Type", "application/json")

	rw.WriteHeader(http.StatusCreated)
	resBytes, _ := json.Marshal(equipmentModel)
	sendJSON(rw, resBytes)

	return
}

func (h *Handler) GetAllEquipmentModels(rw http.ResponseWriter, r *http.Request) {

	const place = "GetAllEquipmentModels"

	var (
		err             error
		equipmentModels []models.EquipmentModel
	)

	rw.Header().Set("Content-Type", "application/json")

	if equipmentModels, err = h.rep.GetAllEquipmentModels(); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		utils.PrintDebug("Error in GetAllEquipmentModels")
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(equipmentModels) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(equipmentModels)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusCreated, place)
	return
}

func (h *Handler) GetEquipmentModelByEquipmentModelUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "GetEquipmentModelByUUID"

	var (
		err                     error
		equipmentModel_uuid     uuid.UUID
		equipmentModel          models.EquipmentModel
		checkFindEquipmentModel bool
	)

	if equipmentModel_uuid, err = getUUID(r, "equipmentModelUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if equipmentModel, checkFindEquipmentModel, err = h.rep.GetEquipmentModelByEquipmentModelUUID(equipmentModel_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	if checkFindEquipmentModel {
		rw.WriteHeader(http.StatusOK)
		resBytes, _ := json.Marshal(equipmentModel)
		sendJSON(rw, resBytes)
	} else {
		sendMessage(rw, "Can't find equipmentModel by uuid: "+equipmentModel_uuid.String(), http.StatusNotFound, nil)
	}

	printResult(err, http.StatusOK, place)
	return
}

func (h *Handler) DelEquipmentModelByEquipmentModelUUID(rw http.ResponseWriter, r *http.Request) {
	const place = "DelEquipmentModelByEquipmentModelUUID"

	var (
		err                     error
		equipmentModel_uuid     uuid.UUID
		checkFindEquipmentModel bool
	)

	if equipmentModel_uuid, err = getUUID(r, "equipmentModelUUID"); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		utils.PrintDebug("Error in getUUID " + err.Error())
		return
	}

	if _, checkFindEquipmentModel, err = h.rep.GetEquipmentModelByEquipmentModelUUID(equipmentModel_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	if !checkFindEquipmentModel {
		sendMessage(rw, "Can't find equipmentModel by uuid: "+equipmentModel_uuid.String(), http.StatusNotFound, nil)
		return
	}

	if err = h.rep.DelEquipmentModelByEquipmentModelUUID(equipmentModel_uuid); err != nil {
		rw.WriteHeader(http.StatusNotFound)
		printResult(err, http.StatusNotFound, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")

	sendMessage(rw, "EquipmentModel was deleted", http.StatusOK, nil)
	printResult(err, http.StatusOK, place)
	return
}
