package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"restService/internal/models"
	"restService/internal/utils"
	"strconv"

	"github.com/gorilla/mux"
)

func (h *Handler) getID(r *http.Request) (personID int, err error) {
	var (
		personIDstr string
		vars        map[string]string
	)

	vars = mux.Vars(r)

	if personIDstr = vars["personID"]; personIDstr == "" {
		err = errors.New("Cant found parameters")
		return
	}

	personID, err = strconv.Atoi(personIDstr)
	if err != nil {
		utils.PrintDebug("strconv error in getID")
		return
	}

	return
}

func getPersonFromBody(r *http.Request) (person models.Person, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of person")
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	_ = json.Unmarshal(body, &person)
	if len(person.Name) == 0 {
		err = errors.New("Cant found name of person")
	}
	return
}
