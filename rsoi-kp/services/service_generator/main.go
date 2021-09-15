package main

import (
	"database/sql"
	"services/config"
	"services/internal/database"
	"services/internal/repInterface"
	api "services/internal/services"
	utils "services/utils"

	"net/http"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	const confPath = "conf.json"

	var (
		db          *sql.DB
		conf        *config.Configuration
		redisClient *redis.Client
		err         error
	)

	if db, redisClient, conf, err = database.DataBaseInitialize(confPath); err != nil {
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
	API := api.NewHandler(pu, redisClient)

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1/generator").Subrouter()
	api.Handle("/start", API.TokenAuthMiddleware(http.HandlerFunc(API.GeneratorStartHandler))).Methods("POST")
	api.Handle("/stop", API.TokenAuthMiddleware(http.HandlerFunc(API.GeneratorStopHandler))).Methods("POST")
	api.Handle("/equipment/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetDataByEquipmentUUID))).Methods("GET")
	api.Handle("/equipment/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DeleteDataByEquipmentUUID))).Methods("DELETE")
	api.Handle("/clear", API.TokenAuthMiddleware(http.HandlerFunc(API.ClearAllData))).Methods("DELETE")
	api.HandleFunc("/service/register", API.ServiceRegister).Methods("POST")

	address := ":" + conf.ServiceGenerator.Port
	utils.PrintDebug("Launched on " + address)

	API.GeneratorStart()
	if err = http.ListenAndServe(address, r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
