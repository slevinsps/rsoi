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
	api := r.PathPrefix("/api/v1/store").Subrouter()
	api.HandleFunc("/{userUid}/purchase", API.PurchaseItem).Methods("POST")
	api.HandleFunc("/{userUid}/orders", API.GetUserOrders).Methods("GET")
	api.HandleFunc("/{userUid}/{orderUid}", API.GetUserOrder).Methods("GET")
	api.HandleFunc("/{userUid}/{orderUid}/warranty", API.RequestWarranty).Methods("POST")
	api.HandleFunc("/{userUid}/{orderUid}/refund", API.ReturnOrder).Methods("DELETE")

	utils.PrintDebug("Launched on " + conf.ServiceStore.Host + ":" + conf.ServiceStore.Port)

	if err = http.ListenAndServe(":"+os.Getenv("PORT"), r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
