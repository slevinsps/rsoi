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
	api := r.PathPrefix("/api/v1/documentation").Subrouter()
	api.Handle("/upload/{equipmentModelUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.CreateFile))).Methods("POST")
	api.Handle("/update", API.TokenAuthMiddleware(http.HandlerFunc(API.UpdateFile))).Methods("PUT")
	api.Handle("/list", API.TokenAuthMiddleware(http.HandlerFunc(API.GetAllFiles))).Methods("GET")
	api.Handle("/equipment_model/{equipmentModelUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetAllFilesByEquipmentModelUUID))).Methods("GET")
	api.Handle("/{fileUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetFileByFileUUID))).Methods("GET")
	api.Handle("/{fileUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DelFileByFileUUID))).Methods("DELETE")
	api.HandleFunc("/service/register", API.ServiceRegister).Methods("POST")

	address := ":" + conf.ServiceDocumenatation.Port
	utils.PrintDebug("Launched on " + address)

	if err = http.ListenAndServe(address, r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
