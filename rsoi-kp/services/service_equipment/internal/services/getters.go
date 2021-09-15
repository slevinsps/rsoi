package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"services/internal/models"
	"services/utils"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func getUUID(r *http.Request, name string) (UID uuid.UUID, err error) {
	var (
		UIDstr string
		vars   map[string]string
	)

	vars = mux.Vars(r)

	if UIDstr = vars[name]; UIDstr == "" {
		err = errors.New("Cant found parameters")
		return
	}

	UID, err = uuid.FromString(UIDstr)
	if err != nil {
		utils.PrintDebug("strconv error in UID")
		return
	}
	return
}

func getEquipmentFromBody(r *http.Request) (equipment models.Equipment, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of equipment")
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	_ = json.Unmarshal(body, &equipment)

	return
}

func getEquipmentModelFromBody(r *http.Request) (equipmentModel models.EquipmentModel, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of equipmentModel")
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	_ = json.Unmarshal(body, &equipmentModel)

	return
}

func getServiceFromBody(r *http.Request) (service models.Service, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of service")
		return
	}
	defer r.Body.Close()

	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	err = json.Unmarshal(body, &service)
	return
}
