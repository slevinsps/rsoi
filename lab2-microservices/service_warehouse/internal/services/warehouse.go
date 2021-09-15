package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"services/config"
	"services/internal/models"
	"services/internal/repInterface"

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

// TakeItem
func (h *Handler) TakeItem(rw http.ResponseWriter, r *http.Request) {
	const place = "TakeItem"

	var (
		err                  error
		orderItem            models.OrderItem
		item                 models.Item
		orderItemTableRecord models.OrderItemTableRecord
	)
	fmt.Println("Warehouse ")
	if orderItem, err = getOrderItemParamsFromBody(r); err != nil {
		sendMessage(rw, "Error in get orderItem", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}
	fmt.Println("Warehouse = ", orderItem)
	if item, err = h.rep.GetItemInfoByModeSize(orderItem.Model, orderItem.Size); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found item", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return
	}

	if item.AvailableCount < 1 {
		sendMessage(rw, "Item not available", http.StatusConflict, nil)
		return
	}

	orderItem.OrderItemUid = uuid.NewV4()
	orderItemTableRecord.Canceled = false
	orderItemTableRecord.ItemID = item.ID
	orderItemTableRecord.OrderUid = orderItem.OrderUid
	orderItemTableRecord.OrderItemUid = orderItem.OrderItemUid

	if err = h.rep.InsertOrderItem(orderItemTableRecord); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	if err = h.rep.TakeOneItem(item); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(http.StatusOK)
	resBytes, _ := json.Marshal(orderItem)
	sendJSON(rw, resBytes)

	printResult(err, http.StatusOK, place)
	return
}

// GetItemInfo
func (h *Handler) GetItemInfo(rw http.ResponseWriter, r *http.Request) {

	const place = "GetItemInfo"

	var (
		err          error
		orderItemUID uuid.UUID
		item         models.Item
	)

	rw.Header().Set("Content-Type", "application/json")

	if orderItemUID, err = getOrderItemUid(r); err != nil {
		sendMessage(rw, "Error in get orderItemUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if item, err = h.rep.GetItemInfoByOrderItemUID(orderItemUID, false); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found item", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return
	}

	rw.WriteHeader(http.StatusOK)

	resBytes, _ := json.Marshal(item)
	sendJSON(rw, resBytes)

	printResult(err, http.StatusOK, place)
	return
}

// ReturnItem
func (h *Handler) ReturnItem(rw http.ResponseWriter, r *http.Request) {
	const place = "ReturnItem"

	var (
		err          error
		orderItemUID uuid.UUID
	)

	rw.Header().Set("Content-Type", "application/json")

	if orderItemUID, err = getOrderItemUid(r); err != nil {
		sendMessage(rw, "Error in get orderItemUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if err = h.rep.ReturnOneItem(orderItemUID); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	if err = h.rep.CancelOrderItem(orderItemUID); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(http.StatusNoContent)

	printResult(err, http.StatusNoContent, place)
	return
}

// WarrantyRequest
func (h *Handler) WarrantyRequest(rw http.ResponseWriter, r *http.Request) {
	const place = "WarrantyRequest"

	var (
		err            error
		orderItemUID   uuid.UUID
		warrantyParams models.WarrantyParams
		item           models.Item
		code           int
		body           []byte
	)

	rw.Header().Set("Content-Type", "application/json")
	if orderItemUID, err = getOrderItemUid(r); err != nil {
		sendMessage(rw, "Error in get orderItemUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if warrantyParams, err = getWarrantyParamsFromBody(r); err != nil {
		sendMessage(rw, "Error in get orderItem", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if item, err = h.rep.GetItemInfoByOrderItemUID(orderItemUID, true); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found item", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return
	}

	warrantyParams.AvailableCount = item.AvailableCount

	if code, body, err = h.WarrantySendRequest(warrantyParams, orderItemUID); err != nil {
		sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusUnprocessableEntity, place)
		return
	}

	if code != http.StatusOK {
		if code == http.StatusNotFound {
			sendMessage(rw, "Warranty not found for itemUid '"+orderItemUID.String()+"'", code, err)
		} else {
			sendMessage(rw, "Error in warranty request", code, err)
		}
		printResult(err, code, place)
		return
	}

	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(code)
	sendJSON(rw, body)

	printResult(err, code, place)
	return
}
