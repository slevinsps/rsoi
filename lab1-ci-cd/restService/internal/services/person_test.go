package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	mock_repInterface "restService/internal/mock"
	"restService/internal/models"
	"strings"
	"testing"

	"io/ioutil"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func Test_PersonCreate(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock_repInterface.NewMockRepInterface(ctl)

	type args struct {
		r        *http.Request
		status   int
		expected models.Person
		times    int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Create person",
			args: args{
				r: httptest.NewRequest("POST", "/persons",
					strings.NewReader(`{"name": "ivan" }`)),
				expected: models.Person{Name: "ivan"},
				status:   http.StatusCreated,
				times:    1,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				personDatabase: mockUsecase,
			}
			w := httptest.NewRecorder()

			mockUsecase.EXPECT().PersonCreate(tt.args.expected).Times(tt.args.times)

			h.PersonCreate(w, tt.args.r)

			if tt.args.status != w.Code {
				t.Error(tt.name)
				t.Error("Expext:", tt.args.status, "and Get:", w.Code)

			}

		})
	}
}

func Test_GetAllPersonsInfo(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock_repInterface.NewMockRepInterface(ctl)

	type args struct {
		r            *http.Request
		status       int
		expected     []models.Person
		expectedBody string
		times        int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Get two person",
			args: args{
				r:            httptest.NewRequest("GET", "/persons", nil),
				expected:     []models.Person{{Name: "ivan", Work: "ya"}, {Name: "vlad", Work: "ozon"}},
				expectedBody: `[{"name":"ivan","work":"ya"},{"name":"vlad","work":"ozon"}]`,
				status:       http.StatusOK,
				times:        1,
			}},
		{
			name: "No person",
			args: args{
				r:            httptest.NewRequest("GET", "/persons", nil),
				expected:     []models.Person{},
				expectedBody: `[]`,
				status:       http.StatusOK,
				times:        1,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				personDatabase: mockUsecase,
			}
			w := httptest.NewRecorder()

			mockUsecase.EXPECT().GetAllPersonsInfo().Return(tt.args.expected, nil).Times(tt.args.times)

			h.GetAllPersonsInfo(w, tt.args.r)

			if tt.args.status != w.Code {
				t.Error(tt.name)
				t.Error("Expext:", tt.args.status, "and Get:", w.Code)
			}

			body, _ := ioutil.ReadAll(w.Body)
			bodyString := string(body)

			if tt.args.expectedBody != bodyString {
				t.Error(tt.name)
				t.Error("Expext:", tt.args.expectedBody, "and Get:", bodyString)
			}
		})
	}
}

func Test_GetPersonByID(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock_repInterface.NewMockRepInterface(ctl)

	type args struct {
		r               *http.Request
		argsID          int
		returnPerson    models.Person
		returnCheckFind bool
		expectedBody    string
		times           int
		status          int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Get exist person",
			args: args{
				r:               httptest.NewRequest("GET", "/persons/1", nil),
				argsID:          1,
				returnPerson:    models.Person{ID: 1, Name: "ivan", Work: "ya"},
				returnCheckFind: true,
				expectedBody:    `{"id":1,"name":"ivan","work":"ya"}`,
				status:          http.StatusOK,
				times:           1,
			}},
		{
			name: "No person",
			args: args{
				r:               httptest.NewRequest("GET", "/persons/100", nil),
				argsID:          100,
				returnCheckFind: false,
				expectedBody:    `{"message":"Can't find person"}`,
				status:          http.StatusNotFound,
				times:           1,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				personDatabase: mockUsecase,
			}
			w := httptest.NewRecorder()

			mockUsecase.EXPECT().GetPersonByID(tt.args.argsID).Return(tt.args.returnPerson, tt.args.returnCheckFind, nil).Times(tt.args.times)

			vars := map[string]string{
				"personID": fmt.Sprint(tt.args.argsID),
			}
			tt.args.r = mux.SetURLVars(tt.args.r, vars)

			h.GetPersonInfo(w, tt.args.r)

			if tt.args.status != w.Code {
				t.Error(tt.name)
				t.Error("Expext:", tt.args.status, "and Get:", w.Code)
			}

			body, _ := ioutil.ReadAll(w.Body)
			bodyString := string(body)

			if tt.args.expectedBody != bodyString {
				t.Error(tt.name)
				t.Error("Expext:", tt.args.expectedBody, "and Get:", bodyString)
			}
		})
	}
}

func Test_UpdatePersonInfo(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock_repInterface.NewMockRepInterface(ctl)

	type args struct {
		r               *http.Request
		argsID          int
		updatedPerson   models.Person
		returnCheckFind bool
		expectedBody    string
		times           int
		status          int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Update exist person",
			args: args{
				r:               httptest.NewRequest("PATCH", "/persons/1", strings.NewReader(`{"name":"ivan","work":"ya"}`)),
				argsID:          1,
				updatedPerson:   models.Person{ID: 1, Name: "ivan", Work: "ya"},
				returnCheckFind: true,
				expectedBody:    `{"id":1,"name":"ivan","work":"ya"}`,
				status:          http.StatusOK,
				times:           1,
			}},
		{
			name: "No person",
			args: args{
				r:               httptest.NewRequest("PATCH", "/persons/100", strings.NewReader(`{"name":"ivan","work":"ya"}`)),
				argsID:          100,
				returnCheckFind: false,
				expectedBody:    `{"message":"Can't find person"}`,
				status:          http.StatusNotFound,
				times:           1,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				personDatabase: mockUsecase,
			}
			w := httptest.NewRecorder()

			mockUsecase.EXPECT().GetPersonByID(tt.args.argsID).Return(tt.args.updatedPerson, tt.args.returnCheckFind, nil).Times(tt.args.times)
			if tt.args.returnCheckFind {
				mockUsecase.EXPECT().UpdatePersonInfo(tt.args.updatedPerson).Return(tt.args.updatedPerson, nil).Times(tt.args.times)
			}

			vars := map[string]string{
				"personID": fmt.Sprint(tt.args.argsID),
			}
			tt.args.r = mux.SetURLVars(tt.args.r, vars)

			h.UpdatePersonInfo(w, tt.args.r)

			if tt.args.status != w.Code {
				t.Error(tt.name)
				t.Error("Expext:", tt.args.status, "and Get:", w.Code)
			}

			body, _ := ioutil.ReadAll(w.Body)
			bodyString := string(body)

			if tt.args.expectedBody != bodyString {
				t.Error(tt.name)
				t.Error("Expext:", tt.args.expectedBody, "and Get:", bodyString)
			}
		})
	}
}

func Test_DeletePersonInfo(t *testing.T) {

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	mockUsecase := mock_repInterface.NewMockRepInterface(ctl)

	type args struct {
		r               *http.Request
		argsID          int
		deletedPerson   models.Person
		returnCheckFind bool
		expectedBody    string
		times           int
		status          int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "Delete exist person",
			args: args{
				r:               httptest.NewRequest("DELETE", "/persons/1", nil),
				argsID:          1,
				deletedPerson:   models.Person{ID: 1, Name: "ivan", Work: "ya"},
				returnCheckFind: true,
				expectedBody:    ``,
				status:          http.StatusOK,
				times:           1,
			}},
		{
			name: "No person",
			args: args{
				r:               httptest.NewRequest("DELETE", "/persons/100", nil),
				argsID:          100,
				returnCheckFind: false,
				expectedBody:    `{"message":"Can't find person"}`,
				status:          http.StatusNotFound,
				times:           1,
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				personDatabase: mockUsecase,
			}
			w := httptest.NewRecorder()

			mockUsecase.EXPECT().GetPersonByID(tt.args.argsID).Return(tt.args.deletedPerson, tt.args.returnCheckFind, nil).Times(tt.args.times)
			if tt.args.returnCheckFind {
				mockUsecase.EXPECT().DeletePersonInfo(tt.args.argsID).Return(nil).Times(tt.args.times)
			}

			vars := map[string]string{
				"personID": fmt.Sprint(tt.args.argsID),
			}
			tt.args.r = mux.SetURLVars(tt.args.r, vars)

			h.DeletePersonInfo(w, tt.args.r)

			if tt.args.status != w.Code {
				t.Error(tt.name)
				t.Error("Expext:", tt.args.status, "and Get:", w.Code)
			}

			body, _ := ioutil.ReadAll(w.Body)
			bodyString := string(body)

			if tt.args.expectedBody != bodyString {
				t.Error(tt.name)
				t.Error("Expext:", tt.args.expectedBody, "and Get:", bodyString)
			}
		})
	}
}
