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
	api := r.PathPrefix("/api/v1/warehouse").Subrouter()
	api.HandleFunc("/{orderItemUid}", API.GetItemInfo).Methods("GET")
	api.HandleFunc("/{orderItemUid}/warranty", API.WarrantyRequest).Methods("POST")
	api.HandleFunc("/", API.TakeItem).Methods("POST")
	api.HandleFunc("/{orderItemUid}", API.ReturnItem).Methods("DELETE")

	utils.PrintDebug("Launched warehouse service on " + conf.ServiceWarehouse.Host + ":" + conf.ServiceWarehouse.Port)

	if err = http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
