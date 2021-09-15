package main

import (
	"database/sql"
	"os"
	"services/config"
	"services/internal/database"
	"services/internal/repInterface"
	api "services/internal/services"
	utils "services/utils"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	const confPath = "conf.json"

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
	api := r.PathPrefix("/api/v1/warranty").Subrouter()
	api.HandleFunc("/{itemUid}", API.StartWarrantyPeriod).Methods("POST")
	api.HandleFunc("/{itemUid}/warranty", API.WarrantyRequest).Methods("POST")
	api.HandleFunc("/{itemUid}", API.GetWarranty).Methods("GET")
	api.HandleFunc("/{itemUid}", API.CloseWarranty).Methods("DELETE")

	utils.PrintDebug("Launched warranty service on " + conf.ServiceWarranty.Host + ":" + conf.ServiceWarranty.Port)

	if err = http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
