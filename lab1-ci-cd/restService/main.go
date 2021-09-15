package main

import (
	"database/sql"
	"os"
	"restService/internal/config"
	"restService/internal/database"
	"restService/internal/repInterface"
	api "restService/internal/services"
	utils "restService/internal/utils"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	const confPath = "restService/conf.json"

	var (
		db   *sql.DB
		conf *config.Configuration
		err  error
	)

	if db, conf, err = database.DataBaseInitialize(confPath); err != nil {
		utils.PrintDebug("Error in configuration file or database " + err.Error())
		return
	}

	Db, err := database.NewDataBase(db)
	if err != nil {
		utils.PrintDebug("Error in database create tabel" + err.Error())
		return
	}
	pu := repInterface.NewPUsecase(Db)
	API := api.NewHandler(pu)

	r := mux.NewRouter()
	r.HandleFunc("/persons", API.PersonCreate).Methods("POST")
	r.HandleFunc("/persons", API.GetAllPersonsInfo).Methods("GET")
	r.HandleFunc("/persons/{personID}", API.GetPersonInfo).Methods("GET")
	r.HandleFunc("/persons/{personID}", API.UpdatePersonInfo).Methods("PATCH")
	r.HandleFunc("/persons/{personID}", API.DeletePersonInfo).Methods("DELETE")

	utils.PrintDebug("Launched on " + conf.Server.Host + ":" + conf.Server.Port)

	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", conf.Server.Port)
	}

	if err = http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
