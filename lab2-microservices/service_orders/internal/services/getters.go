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

func getUserUID(r *http.Request) (userUID uuid.UUID, err error) {
	var (
		userUIDstr string
		vars       map[string]string
	)

	vars = mux.Vars(r)

	if userUIDstr = vars["userUid"]; userUIDstr == "" {
		err = errors.New("Cant found parameters")
		return
	}

	userUID, err = uuid.FromString(userUIDstr)
	if err != nil {
		utils.PrintDebug("strconv error in getUserUID")
		return
	}

	return
}

func getOrderUID(r *http.Request) (orderUID uuid.UUID, err error) {
	var (
		orderUIDstr string
		vars        map[string]string
	)

	vars = mux.Vars(r)

	if orderUIDstr = vars["orderUid"]; orderUIDstr == "" {
		err = errors.New("Cant found parameters")
		return
	}

	orderUID, err = uuid.FromString(orderUIDstr)
	if err != nil {
		utils.PrintDebug("strconv error in getOrderUID")
		return
	}

	return
}

func getItemFromBody(r *http.Request) (item models.Item, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of item")
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	_ = json.Unmarshal(body, &item)
	fmt.Println(item)
	if len(item.Model) == 0 {
		err = errors.New("Cant found model of item")
	}
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
