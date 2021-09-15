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
	api := r.PathPrefix("/api/v1/orders").Subrouter()
	api.HandleFunc("/{userUid}", API.CreateOrder).Methods("POST")
	api.HandleFunc("/{userUid}/{orderUid}", API.GetUserOrderInfo).Methods("GET")
	api.HandleFunc("/{userUid}", API.GetUserOrdersInfo).Methods("GET")
	api.HandleFunc("/{orderUid}/warranty", API.WarrantyRequest).Methods("POST")
	api.HandleFunc("/{orderUid}", API.ReturnOrder).Methods("DELETE")

	utils.PrintDebug("Launched service orders on " + conf.ServiceOrders.Host + ":" + conf.ServiceOrders.Port)

	if err = http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
