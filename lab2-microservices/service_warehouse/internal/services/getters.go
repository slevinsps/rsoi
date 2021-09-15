package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"services/internal/models"
	"services/utils"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func getOrderItemUid(r *http.Request) (orderItemUID uuid.UUID, err error) {
	var (
		orderItemUIDstr string
		vars            map[string]string
	)

	vars = mux.Vars(r)

	if orderItemUIDstr = vars["orderItemUid"]; orderItemUIDstr == "" {
		err = errors.New("Cant found parameters")
		return
	}

	orderItemUID, err = uuid.FromString(orderItemUIDstr)
	if err != nil {
		utils.PrintDebug("strconv error in getItemUID")
		return
	}

	return
}

func getOrderItemParamsFromBody(r *http.Request) (orderItem models.OrderItem, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of item")
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	_ = json.Unmarshal(body, &orderItem)
	fmt.Println(orderItem)

	return
}

func getWarrantyParamsFromBody(r *http.Request) (warrantyParams models.WarrantyParams, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of item")
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	_ = json.Unmarshal(body, &warrantyParams)
	fmt.Println(warrantyParams)

	return
}
