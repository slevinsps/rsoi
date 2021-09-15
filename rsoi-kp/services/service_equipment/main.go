package main

import (
	"database/sql"
	"services/config"
	"services/internal/database"
	"services/internal/repInterface"
	api "services/internal/services"
	utils "services/utils"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

	if err = godotenv.Load(); err != nil {
		utils.PrintDebug("Error loading .env file")
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
	api := r.PathPrefix("/api/v1/equipment").Subrouter()
	api.Handle("/create", API.TokenAuthMiddleware(http.HandlerFunc(API.CreateEquipment))).Methods("POST")
	api.Handle("/update", API.TokenAuthMiddleware(http.HandlerFunc(API.UpdateEquipment))).Methods("PUT")
	api.Handle("/list", API.TokenAuthMiddleware(http.HandlerFunc(API.GetAllEquipments))).Methods("GET")
	api.Handle("/list/{monitorUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetEquipmentsByMonitorUUID))).Methods("GET")
	api.Handle("/list/{monitorUUID}/notadded", API.TokenAuthMiddleware(http.HandlerFunc(API.GetNotAddedEquipmentsByMonitorUUID))).Methods("GET")
	api.Handle("/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetEquipmentByEquipmentUUID))).Methods("GET")
	api.Handle("/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DelEquipmentByEquipmentUUID))).Methods("DELETE")

	api.Handle("/model/create", API.TokenAuthMiddleware(http.HandlerFunc(API.CreateEquipmentModel))).Methods("POST")
	api.Handle("/model/list", API.TokenAuthMiddleware(http.HandlerFunc(API.GetAllEquipmentModels))).Methods("GET")
	api.Handle("/model/{equipmentModelUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetEquipmentModelByEquipmentModelUUID))).Methods("GET")
	api.Handle("/model/{equipmentModelUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DelEquipmentModelByEquipmentModelUUID))).Methods("DELETE")
	api.HandleFunc("/service/register", API.ServiceRegister).Methods("POST")

	address := ":" + conf.ServiceEquipment.Port
	utils.PrintDebug("Launched on " + address)

	if err = http.ListenAndServe(address, r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
