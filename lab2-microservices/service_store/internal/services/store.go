package api

import (
	"database/sql"
	"encoding/json"
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

// PurchaseItem
func (h *Handler) PurchaseItem(rw http.ResponseWriter, r *http.Request) {
	const place = "PurchaseItem"

	var (
		err      error
		userUID  uuid.UUID
		code     int
		item     models.Item
		orderUid models.OrderUid
	)

	if userUID, err = getUserUID(r); err != nil {
		sendMessage(rw, "Error in get orderUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if _, err = h.rep.GetUserByUUID(userUID); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found user", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return
	}

	if item, err = getItemFromBody(r); err != nil {
		sendMessage(rw, "Error in get item params", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code, orderUid, err = h.CreateOrder(item, userUID); err != nil {
		sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code != http.StatusOK {
		sendMessage(rw, "Error in external request", code, err)
		printResult(err, code, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Location", "https://store/"+orderUid.OrderUid.String())

	rw.WriteHeader(http.StatusCreated)

	printResult(err, http.StatusCreated, place)
	return
}

// RequestWarranty
func (h *Handler) RequestWarranty(rw http.ResponseWriter, r *http.Request) {
	const place = "RequestWarranty"

	var (
		err              error
		orderUID         uuid.UUID
		userUID          uuid.UUID
		code             int
		warrantyParams   models.WarrantyParams
		warrantyResponse models.WarrantyResponse
		body             []byte
	)

	rw.Header().Set("Content-Type", "application/json")

	if userUID, err = getUserUID(r); err != nil {
		sendMessage(rw, "Error in get orderUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if _, err = h.rep.GetUserByUUID(userUID); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found user", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return
	}

	if orderUID, err = getOrderUID(r); err != nil {
		sendMessage(rw, "Error in get orderUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if warrantyParams, err = getWarrantyParamsFromBody(r); err != nil {
		sendMessage(rw, "Error in get reason from body", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code, warrantyResponse, err = h.RequestWarrantyInOrders(orderUID, warrantyParams); err != nil {
		sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code != http.StatusOK {
		sendMessage(rw, "Error in external request", code, err)
		printResult(err, code, place)
		return
	}

	warrantyResponse.OrderUID = orderUID

	body, _ = json.Marshal(warrantyResponse)
	rw.WriteHeader(http.StatusOK)
	sendJSON(rw, body)

	printResult(err, http.StatusOK, place)
	return
}

// ReturnOrder
func (h *Handler) ReturnOrder(rw http.ResponseWriter, r *http.Request) {
	const place = "ReturnOrder"

	var (
		err      error
		orderUID uuid.UUID
		userUID  uuid.UUID
		code     int
	)

	rw.Header().Set("Content-Type", "application/json")
	if userUID, err = getUserUID(r); err != nil {
		sendMessage(rw, "Error in get orderUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if _, err = h.rep.GetUserByUUID(userUID); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found user", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return
	}

	if orderUID, err = getOrderUID(r); err != nil {
		sendMessage(rw, "Error in get orderUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code, err = h.ReturnOrderInOrders(orderUID); err != nil {
		sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code != http.StatusNoContent {
		sendMessage(rw, "Error in external request", code, err)
		printResult(err, code, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(code)

	printResult(err, code, place)
	return
}

// GetUserOrders
func (h *Handler) GetUserOrders(rw http.ResponseWriter, r *http.Request) {
	const place = "GetUserOrders"

	var (
		err                error
		userUID            uuid.UUID
		code               int
		ordersArray        []models.Order
		orderResponseArray []models.OrderResponse
		orderResponse      models.OrderResponse
		item               models.Item
		warranty           models.Warranty
		body               []byte
		// orderUid models.OrderUid
	)

	if userUID, err = getUserUID(r); err != nil {
		sendMessage(rw, "Error in get orderUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if _, err = h.rep.GetUserByUUID(userUID); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found user", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return
	}

	if code, ordersArray, err = h.UserOrdersInOrders(userUID); err != nil {
		sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code != http.StatusOK {
		sendMessage(rw, "Error in external request", code, err)
		printResult(err, code, place)
		return
	}

	for _, order := range ordersArray {
		orderResponse.OrderUID = order.OrderUID
		orderResponse.Date = order.OrderDate.String()

		if code, item, err = h.GetItemInWarehouse(order.ItemUID); err != nil {
			sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
			printResult(err, http.StatusBadRequest, place)
			return
		}

		if code == http.StatusOK {
			orderResponse.Model = item.Model
			orderResponse.Size = item.Size
		}

		if code, warranty, err = h.CheckWarrantyStatusInWarranty(order.ItemUID); err != nil {
			sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
			printResult(err, http.StatusBadRequest, place)
			return
		}

		if code == http.StatusOK {
			orderResponse.WarrantyDate = warranty.WarrantyDate.String()
			orderResponse.WarrantyStatus = warranty.Status
		}

		orderResponseArray = append(orderResponseArray, orderResponse)
	}

	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(http.StatusOK)
	if len(orderResponseArray) == 0 {
		rw.Write([]byte("[]"))
	} else {
		body, _ = json.Marshal(orderResponseArray)
		sendJSON(rw, body)
	}

	printResult(err, http.StatusOK, place)
	return
}

// GetUserOrder
func (h *Handler) GetUserOrder(rw http.ResponseWriter, r *http.Request) {
	const place = "GetUserOrder"

	var (
		err           error
		userUID       uuid.UUID
		orderUID      uuid.UUID
		code          int
		order         models.Order
		orderResponse models.OrderResponse
		item          models.Item
		warranty      models.Warranty
		body          []byte
		// orderUid models.OrderUid
	)

	if userUID, err = getUserUID(r); err != nil {
		sendMessage(rw, "Error in get userUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if _, err = h.rep.GetUserByUUID(userUID); err != nil {
		if err == sql.ErrNoRows {
			sendMessage(rw, "Not found user", http.StatusNotFound, nil)
		} else {
			sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
			printResult(err, http.StatusInternalServerError, place)
		}
		return
	}

	if orderUID, err = getOrderUID(r); err != nil {
		sendMessage(rw, "Error in get orderUID", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code, order, err = h.UserOrderInOrders(userUID, orderUID); err != nil {
		sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code != http.StatusOK {
		if code == http.StatusNotFound {
			sendMessage(rw, "Order not found", code, err)
			printResult(err, code, place)
		} else {
			sendMessage(rw, "Error in external request", code, err)
			printResult(err, code, place)
		}
		return
	}

	orderResponse.OrderUID = orderUID
	orderResponse.Date = order.OrderDate.String()

	if code, item, err = h.GetItemInWarehouse(order.ItemUID); err != nil {
		sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code == http.StatusOK {
		orderResponse.Model = item.Model
		orderResponse.Size = item.Size
	}

	if code, warranty, err = h.CheckWarrantyStatusInWarranty(order.ItemUID); err != nil {
		sendMessage(rw, "Error in external request", http.StatusUnprocessableEntity, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if code == http.StatusOK {
		orderResponse.WarrantyDate = warranty.WarrantyDate.String()
		orderResponse.WarrantyStatus = warranty.Status
	}

	rw.Header().Set("Content-Type", "application/json")
	// rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	body, _ = json.Marshal(orderResponse)
	rw.WriteHeader(http.StatusOK)
	sendJSON(rw, body)

	printResult(err, http.StatusOK, place)
	return
}
