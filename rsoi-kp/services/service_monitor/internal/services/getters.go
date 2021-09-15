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

func getMonitorFromBody(r *http.Request) (monitor models.Monitor, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of monitor")
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err = json.Unmarshal(body, &monitor)

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
