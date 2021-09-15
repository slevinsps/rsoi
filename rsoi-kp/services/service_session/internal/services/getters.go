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

func getUserLogin(r *http.Request) (login string, err error) {
	var (
		vars map[string]string
	)

	vars = mux.Vars(r)
	login = vars["userUid"]
	if login == "" {
		err = errors.New("Cant found login")
		return
	}

	return
}

func getUserFromBody(r *http.Request) (user models.User, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of user")
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	_ = json.Unmarshal(body, &user)

	fmt.Println(user.Login)
	fmt.Println(user.Password)
	if len(user.Login) == 0 {
		err = errors.New("Login is empty")
	}
	if len(user.Password) == 0 {
		err = errors.New("Password is empty")
	}
	return
}

func getTokenDetailsFromBody(r *http.Request) (td models.TokenDetails, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of user")
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	_ = json.Unmarshal(body, &td)

	if len(td.RefreshToken) == 0 {
		err = errors.New("Refresh token is empty")
	}

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

func getTDFromBody(r *http.Request) (td models.TokenDetails, err error) {
	if r.Body == nil || r.ContentLength == 0 {
		err = errors.New("Cant found parameters of ad")
		return
	}
	defer r.Body.Close()

	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	err = json.Unmarshal(body, &td)
	return
}
