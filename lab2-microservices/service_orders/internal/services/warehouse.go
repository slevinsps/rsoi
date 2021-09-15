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

func (h *Handler) getItemUid(item models.Item, orderUID uuid.UUID) (code int, respItemOrder models.ItemOrder, err error) {
	var (
		resp      *http.Response
		itemOrder models.ItemOrder
		body      []byte
	)

	itemOrder.Model = item.Model
	itemOrder.Size = item.Size
	itemOrder.OrderUid = orderUID

	sendBody, _ := json.Marshal(itemOrder)
	fmt.Println("itemOrder = ", itemOrder)
	urlString := h.conf.ServiceWarehouse.URL + "/api/v1/warehouse/"
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
	err = json.Unmarshal(body, &respItemOrder)

	return
}

func (h *Handler) returnItem(itemUID uuid.UUID) (code int, err error) {
	var (
		resp *http.Response
	)

	urlString := h.conf.ServiceWarehouse.URL + "/api/v1/warehouse/" + itemUID.String()
	fmt.Println(urlString)

	resp, err = utils.Delete(urlString, []byte{})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = resp.StatusCode

	return
}

func (h *Handler) requestWarranty(itemUID uuid.UUID, warrantyParams models.WarrantyParams) (code int, body []byte, err error) {
	var (
		resp *http.Response
	)

	urlString := h.conf.ServiceWarehouse.URL + "/api/v1/warehouse/" + itemUID.String() + "/warranty"
	fmt.Println(urlString)

	sendBody, _ := json.Marshal(warrantyParams)
	resp, err = utils.Post(urlString, sendBody)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	code = http.StatusOK
	body, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		code = http.StatusUnprocessableEntity
	}

	return
}
