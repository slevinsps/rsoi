package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"services/internal/models"
	"services/utils"

	uuid "github.com/satori/go.uuid"
)

// CreateOrder
func (h *Handler) CreateOrder(item models.Item, userUID uuid.UUID) (code int, orderUid models.OrderUid, err error) {
	var (
		resp *http.Response
		body []byte
	)

	sendBody, _ := json.Marshal(item)
	fmt.Println("itemOrder = ", item)
	urlString := h.conf.ServiceOrders.URL + "/api/v1/orders/" + userUID.String()
	fmt.Println(urlString)
	resp, err = utils.Post(urlString, sendBody)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)

	if code != http.StatusOK {
		return
	}
	err = json.Unmarshal(body, &orderUid)

	return
}

// RequestWarranty
func (h *Handler) RequestWarrantyInOrders(orderUID uuid.UUID, reason models.WarrantyParams) (code int, warrantyResponse models.WarrantyResponse, err error) {
	var (
		resp *http.Response
		body []byte
	)

	urlString := h.conf.ServiceOrders.URL + "/api/v1/orders/" + orderUID.String() + "/warranty"
	fmt.Println(urlString)
	sendBody, _ := json.Marshal(reason)
	fmt.Println("reason = ", reason)
	resp, err = utils.Post(urlString, sendBody)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)

	if code != http.StatusOK {
		return
	}
	err = json.Unmarshal(body, &warrantyResponse)

	return
}

// ReturnOrderInOrders
func (h *Handler) ReturnOrderInOrders(orderUID uuid.UUID) (code int, err error) {
	var (
		resp *http.Response
	)

	urlString := h.conf.ServiceOrders.URL + "/api/v1/orders/" + orderUID.String()
	fmt.Println(urlString)
	resp, err = utils.Delete(urlString, []byte{})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode

	return
}

// UserOrdersInOrders
func (h *Handler) UserOrdersInOrders(userUID uuid.UUID) (code int, ordersArray []models.Order, err error) {
	var (
		resp *http.Response
		body []byte
	)

	urlString := h.conf.ServiceOrders.URL + "/api/v1/orders/" + userUID.String()
	fmt.Println(urlString)

	resp, err = utils.Get(urlString, []byte{})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)

	if code != http.StatusOK {
		return
	}
	err = json.Unmarshal(body, &ordersArray)

	return
}

// UserOrderInOrders
func (h *Handler) UserOrderInOrders(userUID uuid.UUID, orderUID uuid.UUID) (code int, order models.Order, err error) {
	var (
		resp *http.Response
		body []byte
	)

	urlString := h.conf.ServiceOrders.URL + "/api/v1/orders/" + userUID.String() + "/" + orderUID.String()
	fmt.Println(urlString)

	resp, err = utils.Get(urlString, []byte{})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode
	body, err = ioutil.ReadAll(resp.Body)

	if code != http.StatusOK {
		return
	}
	err = json.Unmarshal(body, &order)

	return
}
