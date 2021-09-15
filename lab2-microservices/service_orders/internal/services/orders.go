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

// CreateOrder
func (h *Handler) CreateOrder(rw http.ResponseWriter, r *http.Request) {
	const place = "CreateOrder"

	var (
		err           error
		userUID       uuid.UUID
		itemUID       uuid.UUID
		item          models.Item
		respOrderItem models.ItemOrder
		code          int
		resOrder      models.OrderUid
	)

	orderUID := uuid.NewV4()

	if userUID, err = getUserUID(r); err != nil {
		sendMessage(rw, "Error in get useruid", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if item, err = getItemFromBody(r); err != nil {
		sendMessage(rw, "Error in get item from request", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code, respOrderItem, err = h.getItemUid(item, orderUID); err != nil {
		sendMessage(rw, "Error in get item UUID", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusUnprocessableEntity, place)
		return
	}

	if code == http.StatusConflict {
		sendMessage(rw, "Item not available", http.StatusConflict, nil)
		return
	}

	if code != http.StatusOK {
		fmt.Println(code)
		sendMessage(rw, "External request to warehouse failed", http.StatusUnprocessableEntity, nil)
		return
	}

	itemUID = respOrderItem.OrderItemUid

	if err = h.startWarranty(itemUID); err != nil {
		sendMessage(rw, "Error in start warranty", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if err = h.rep.OrderCreate(orderUID, userUID, itemUID); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(http.StatusOK)
	resOrder.OrderUid = orderUID
	resBytes, _ := json.Marshal(resOrder)
	sendJSON(rw, resBytes)

	printResult(err, http.StatusOK, place)
	return
}

// GetUserOrderInfo
func (h *Handler) GetUserOrderInfo(rw http.ResponseWriter, r *http.Request) {

	const place = "GetAllPersonsInfo"

	var (
		err      error
		userUID  uuid.UUID
		orderUID uuid.UUID
		order    models.Order
	)

	rw.Header().Set("Content-Type", "application/json")

	if orderUID, err = getOrderUID(r); err != nil {
		sendMessage(rw, "Error in get orderUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}
	if userUID, err = getUserUID(r); err != nil {
		sendMessage(rw, "Error in get userUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if order, err = h.rep.GetUserOrderInfo(userUID, orderUID); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found item", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return
	}

	rw.WriteHeader(http.StatusOK)

	resBytes, _ := json.Marshal(order)
	sendJSON(rw, resBytes)

	printResult(err, http.StatusOK, place)
	return
}

// GetUserOrdersInfo
func (h *Handler) GetUserOrdersInfo(rw http.ResponseWriter, r *http.Request) {

	const place = "GetUserOrdersInfo"

	var (
		err     error
		userUID uuid.UUID
		orders  []models.Order
	)

	rw.Header().Set("Content-Type", "application/json")

	if userUID, err = getUserUID(r); err != nil {
		sendMessage(rw, "Error in get userUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if orders, err = h.rep.GetUserOrdersInfo(userUID); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(orders) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(orders)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusOK, place)
	return
}

// WarrantyRequest
func (h *Handler) WarrantyRequest(rw http.ResponseWriter, r *http.Request) {
	const place = "WarrantyRequest"

	var (
		err            error
		orderUID       uuid.UUID
		itemUID        uuid.UUID
		order          models.Order
		code           int
		warrantyParams models.WarrantyParams
		body           []byte
	)

	if orderUID, err = getOrderUID(r); err != nil {
		sendMessage(rw, "Error in get orderUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if order, err = h.rep.GetOrderInfoByOrderUID(orderUID); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found order", http.StatusNotFound, nil)
			return
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
			return
		}
	}

	if warrantyParams, err = getWarrantyParamsFromBody(r); err != nil {
		sendMessage(rw, "Error in get warrantyParams", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	itemUID = order.ItemUID

	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	if code, body, err = h.requestWarranty(itemUID, warrantyParams); err != nil {
		sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.WriteHeader(code)
	sendJSON(rw, body)

	printResult(err, code, place)
	return
}

// ReturnOrder
func (h *Handler) ReturnOrder(rw http.ResponseWriter, r *http.Request) {
	const place = "ReturnOrder"

	var (
		err      error
		orderUID uuid.UUID
		itemUID  uuid.UUID
		order    models.Order
		code     int
	)

	if orderUID, err = getOrderUID(r); err != nil {
		sendMessage(rw, "Error in get orderUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if order, err = h.rep.GetOrderInfoByOrderUID(orderUID); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found order", http.StatusNotFound, nil)
			return
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
			return
		}
	}

	itemUID = order.ItemUID

	if code, err = h.returnItem(itemUID); err != nil {
		sendMessage(rw, "Error in returnItem", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusUnprocessableEntity, place)
		return
	}

	if code, err = h.stopWarranty(itemUID); err != nil {
		sendMessage(rw, "Error in stopWarranty", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusUnprocessableEntity, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(code)

	printResult(err, code, place)
	return
}
