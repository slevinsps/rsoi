package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"services/constants"
	"services/internal/models"
	"services/internal/repInterface"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Handler struct {
	rep repInterface.RepInterface
}

func NewHandler(rep repInterface.RepInterface) *Handler {
	return &Handler{rep: rep}
}

func sendMessage(rw http.ResponseWriter, messageText string, status int, err error) {
	rw.WriteHeader(status)
	message := models.Message{Message: messageText}
	if err != nil {
		additionalProps := &models.AdditionalProp{AdditionalProp1: err.Error()}
		message.Errors = additionalProps
	}
	resBytes, _ := json.Marshal(message)
	sendJSON(rw, resBytes)
}

// StartWarrantyPeriod
func (h *Handler) StartWarrantyPeriod(rw http.ResponseWriter, r *http.Request) {
	const place = "StartWarrantyPeriod"

	var (
		err     error
		itemUID uuid.UUID
	)

	if itemUID, err = getItemUID(r); err != nil {
		sendMessage(rw, "Error in get itemUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if err = h.rep.StartWarrantyPeriod(itemUID); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(http.StatusNoContent)

	printResult(err, http.StatusCreated, place)
	return
}

// GetWarranty
func (h *Handler) GetWarranty(rw http.ResponseWriter, r *http.Request) {

	const place = "GetWarranty"

	var (
		err      error
		itemUID  uuid.UUID
		warranty models.Warranty
	)

	rw.Header().Set("Content-Type", "application/json")

	if itemUID, err = getItemUID(r); err != nil {
		sendMessage(rw, "Error in get itemUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if warranty, err = h.rep.GetWarranty(itemUID); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found warranty", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return
	}

	rw.WriteHeader(http.StatusOK)

	resBytes, _ := json.Marshal(warranty)
	sendJSON(rw, resBytes)

	printResult(err, http.StatusOK, place)
	return
}

func checkWarranty(warranty models.Warranty) (checkedWarranty bool) {
	now := time.Now()
	nowMinusMonth := now.AddDate(0, -1, 0)
	fmt.Println("nowMinusMonth ", nowMinusMonth)
	fmt.Println("warranty.WarrantyDate ", warranty.WarrantyDate)
	checkedWarranty = warranty.WarrantyDate.After(nowMinusMonth)
	return
}

// WarrantyRequest
func (h *Handler) WarrantyRequest(rw http.ResponseWriter, r *http.Request) {
	const place = "WarrantyRequest"

	var (
		err              error
		itemUID          uuid.UUID
		warranty         models.Warranty
		decision         string
		warrantyParams   models.WarrantyParams
		waarantyResponse models.WarrantyResponse
	)

	if itemUID, err = getItemUID(r); err != nil {
		sendMessage(rw, "Error in get itemUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if warrantyParams, err = getWarrantyParamsFromBody(r); err != nil {
		sendMessage(rw, "Error in get warranty params", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if warranty, err = h.rep.GetWarranty(itemUID); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found warranty", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return

	}

	decision = constants.WARRANTY_DECISION_RESUSE
	if warranty.Status == constants.WARRANTY_STATUS_ON && checkWarranty(warranty) {
		if warrantyParams.AvailableCount > 0 {
			decision = constants.WARRANTY_DECISION_RETURN
		} else {
			decision = constants.WARRANTY_DECISION_FIXING
		}
	}

	fmt.Println("decision = ", decision)
	fmt.Println("warranty.Status = ", warranty.Status)

	warranty.Comment = warrantyParams.Reason
	if decision == constants.WARRANTY_DECISION_RESUSE {
		warranty.Status = constants.WARRANTY_STATUS_REMOVED
	} else {
		warranty.Status = constants.WARRANTY_STATUS_USE
	}

	fmt.Println("warranty.Status = ", warranty.Status)
	fmt.Println("warranty.Comment = ", warranty.Comment)

	if err = h.rep.UpdateWarranty(warranty); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(http.StatusOK)
	waarantyResponse.WarrantyDate = warranty.WarrantyDate.String()
	waarantyResponse.Decision = decision
	resBytes, _ := json.Marshal(waarantyResponse)
	sendJSON(rw, resBytes)

	printResult(err, http.StatusOK, place)
	return
}

// CloseWarranty
func (h *Handler) CloseWarranty(rw http.ResponseWriter, r *http.Request) {
	const place = "CloseWarranty"

	var (
		err     error
		itemUID uuid.UUID
	)

	if itemUID, err = getItemUID(r); err != nil {
		sendMessage(rw, "Error in get itemUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if err = h.rep.CloseWarranty(itemUID); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(http.StatusNoContent)

	printResult(err, http.StatusNoContent, place)
	return
}
