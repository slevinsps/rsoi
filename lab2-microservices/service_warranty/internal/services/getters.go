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

func getItemUID(r *http.Request) (itemUID uuid.UUID, err error) {
	var (
		itemUIDstr string
		vars       map[string]string
	)

	vars = mux.Vars(r)

	if itemUIDstr = vars["itemUid"]; itemUIDstr == "" {
		err = errors.New("Cant found parameters")
		return
	}

	itemUID, err = uuid.FromString(itemUIDstr)
	if err != nil {
		utils.PrintDebug("strconv error in getItemUID")
		return
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
	if len(warrantyParams.Reason) == 0 {
		err = errors.New("Cant found model of item")
	}
	return
}
