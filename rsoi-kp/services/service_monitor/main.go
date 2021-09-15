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
	route := r.PathPrefix("/api/v1/monitor").Subrouter()
	route.Handle("/create", API.TokenAuthMiddleware(http.HandlerFunc(API.CreateMonitor))).Methods("POST")
	route.Handle("/{monitorUUID}/add/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.AddEquipment))).Methods("POST")
	route.Handle("/{monitorUUID}/del/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DelEquipment))).Methods("DELETE")

	route.Handle("/user/{userUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetAllMonitorsByUserUUID))).Methods("GET")
	route.Handle("/{monitorUUID}/user/{userUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetMonitorByMonitorUUIDuserUUID))).Methods("GET")
	route.Handle("/{monitorUUID}/user/{userUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DelMonitorByMonitorUUIDuserUUID))).Methods("DELETE")
	route.HandleFunc("/service/register", API.ServiceRegister).Methods("POST")

	address := ":" + conf.ServiceMonitor.Port
	utils.PrintDebug("Launched on " + address)

	if err = http.ListenAndServe(address, r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
