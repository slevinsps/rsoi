package api

import (
	"encoding/json"
	"net/http"
	"restService/internal/models"
	"restService/internal/repInterface"
	"strconv"
)

type Handler struct {
	personDatabase repInterface.RepInterface
}

func NewHandler(personDatabase repInterface.RepInterface) *Handler {
	return &Handler{personDatabase: personDatabase}
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

func (h *Handler) PersonCreate(rw http.ResponseWriter, r *http.Request) {
	const place = "PersonCreate"

	var (
		err           error
		person        models.Person
		createdPerson models.Person
	)

	if person, err = getPersonFromBody(r); err != nil {
		sendMessage(rw, "Error in get person from request", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if createdPerson, err = h.personDatabase.PersonCreate(person); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Location", "https://rsoi-person-service.herokuapp.com/persons/"+strconv.Itoa(createdPerson.ID))

	rw.WriteHeader(http.StatusCreated)

	printResult(err, http.StatusCreated, place)
	return
}

// GetAllPersonsInfo
func (h *Handler) GetAllPersonsInfo(rw http.ResponseWriter, r *http.Request) {

	const place = "GetAllPersonsInfo"

	var (
		err     error
		persons []models.Person
	)

	rw.Header().Set("Content-Type", "application/json")

	if persons, err = h.personDatabase.GetAllPersonsInfo(); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if len(persons) == 0 {
		rw.Write([]byte("[]"))
	} else {
		resBytes, _ := json.Marshal(persons)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusCreated, place)
	return
}

// GetPersonInfo
func (h *Handler) GetPersonInfo(rw http.ResponseWriter, r *http.Request) {

	const place = "GetPersonInfo"

	var (
		err             error
		ID              int
		checkFindPerson bool
		person          models.Person
	)

	if ID, err = h.getID(r); err != nil {
		sendMessage(rw, "Error in get person id", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if person, checkFindPerson, err = h.personDatabase.GetPersonByID(ID); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	if !checkFindPerson {
		sendMessage(rw, "Can't find person", http.StatusNotFound, nil)
	} else {
		rw.WriteHeader(http.StatusOK)
		resBytes, _ := json.Marshal(person)
		sendJSON(rw, resBytes)
	}

	printResult(err, http.StatusCreated, place)
	return
}

// UpdatePersonInfo
func (h *Handler) UpdatePersonInfo(rw http.ResponseWriter, r *http.Request) {

	const place = "UpdatePersonInfo"

	var (
		err             error
		ID              int
		newPerson       models.Person
		checkFindPerson bool
	)
	rw.Header().Set("Content-Type", "application/json")

	if ID, err = h.getID(r); err != nil {
		sendMessage(rw, "Error in get person id", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if newPerson, err = getPersonFromBody(r); err != nil {
		sendMessage(rw, "Error in get person from request", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if _, checkFindPerson, err = h.personDatabase.GetPersonByID(ID); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}

	if !checkFindPerson {
		sendMessage(rw, "Can't find person", http.StatusNotFound, nil)
		return
	}

	newPerson.ID = ID

	if _, err = h.personDatabase.UpdatePersonInfo(newPerson); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}
	rw.WriteHeader(http.StatusOK)
	resBytes, _ := json.Marshal(newPerson)
	sendJSON(rw, resBytes)

	printResult(err, http.StatusOK, place)
	return
}

// DeletePersonInfo
func (h *Handler) DeletePersonInfo(rw http.ResponseWriter, r *http.Request) {

	const place = "DeletePersonInfo"

	var (
		err             error
		ID              int
		checkFindPerson bool
	)
	rw.Header().Set("Content-Type", "application/json")

	if ID, err = h.getID(r); err != nil {
		sendMessage(rw, "Error in get person id", http.StatusBadRequest, err)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if _, checkFindPerson, err = h.personDatabase.GetPersonByID(ID); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusBadRequest, place)
		return
	}

	if !checkFindPerson {
		sendMessage(rw, "Can't find person", http.StatusNotFound, nil)
		return
	}

	if err = h.personDatabase.DeletePersonInfo(ID); err != nil {
		sendMessage(rw, "Error in database", http.StatusInternalServerError, nil)
		printResult(err, http.StatusInternalServerError, place)
		return
	}
	rw.WriteHeader(http.StatusOK)
	printResult(err, http.StatusOK, place)
	return
}
