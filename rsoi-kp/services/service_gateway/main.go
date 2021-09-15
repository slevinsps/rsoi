package main

import (
	"database/sql"
	"fmt"
	"os"
	"services/config"
	"services/internal/database"
	"services/internal/models"
	api "services/internal/services"
	utils "services/utils"

	"net/http"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

func addServicesSecretsToDB(db *database.DataBase) (err error) {
	var (
		service models.Service
	)
	services := [2]string{"GENERATOR", "GATEWAY"}
	for i := 0; i < len(services); i++ {
		service.Login = os.Getenv(services[i] + "_ID")
		service.Password = os.Getenv(services[i] + "_SECRET")
		service.Password, err = api.HashPassword(service.Password)
		if err = db.AddServiceSecrets(service); err != nil {
			return
		}
	}
	return
}

func main() {
	const confPath = "conf.json"

	var (
		conf        *config.Configuration
		redisClient *redis.Client
		err         error
		db          *sql.DB
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
	API := api.NewHandler(Db, redisClient)

	r := mux.NewRouter()
	route := r.PathPrefix("/api/v1/app/").Subrouter()

	route.Handle("/monitor/create", API.TokenAuthMiddleware(http.HandlerFunc(API.CreateMonitor), false)).Methods("POST")
	route.Handle("/monitor/{monitorUUID}/add/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.AddEquipment), false)).Methods("POST")
	route.Handle("/monitor/{monitorUUID}/del/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DelEquipment), false)).Methods("DELETE")
	route.Handle("/monitor/list", API.TokenAuthMiddleware(http.HandlerFunc(API.GetAllMonitorsByUserUUID), false)).Methods("GET")
	route.Handle("/monitor/{monitorUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetMonitorByMonitorUUID), false)).Methods("GET")
	route.Handle("/monitor/{monitorUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DelMonitorByMonitorUUID), false)).Methods("DELETE")

	route.Handle("/equipment/create", API.TokenAuthMiddleware(http.HandlerFunc(API.CreateEquipment), true)).Methods("POST")
	route.Handle("/equipment/update", API.TokenAuthMiddleware(http.HandlerFunc(API.UpdateEquipment), true)).Methods("PUT")
	route.Handle("/equipment/list", API.TokenAuthMiddleware(http.HandlerFunc(API.GetAllEquipments), true)).Methods("GET")
	route.Handle("/equipment/list/{monitorUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetEquipmentsByMonitorUUID), false)).Methods("GET")
	route.Handle("/equipment/list/{monitorUUID}/notadded", API.TokenAuthMiddleware(http.HandlerFunc(API.GetNotAddedEquipmentsByMonitorUUID), false)).Methods("GET")
	route.Handle("/equipment/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetEquipmentByEquipmentUUID), false)).Methods("GET")
	route.Handle("/equipment/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DelEquipmentByEquipmentUUID), true)).Methods("DELETE")
	route.Handle("/equipment/model/create", API.TokenAuthMiddleware(http.HandlerFunc(API.CreateEquipmentModel), true)).Methods("POST")
	route.Handle("/equipment/model/list", API.TokenAuthMiddleware(http.HandlerFunc(API.GetAllEquipmentModels), false)).Methods("GET")
	route.Handle("/equipment/model/{equipmentModelUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetEquipmentModelByEquipmentModelUUID), false)).Methods("GET")
	route.Handle("/equipment/model/{equipmentModelUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DelEquipmentModelByEquipmentModelUUID), true)).Methods("DELETE")

	route.Handle("/documentation/upload/{equipmentModelUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.CreateFile), true)).Methods("POST")
	route.Handle("/documentation/update", API.TokenAuthMiddleware(http.HandlerFunc(API.UpdateFile), true)).Methods("PUT")
	route.Handle("/documentation/list", API.TokenAuthMiddleware(http.HandlerFunc(API.GetAllFiles), false)).Methods("GET")
	route.Handle("/documentation/equipment_model/{equipmentModelUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetAllFilesByEquipmentUUID), false)).Methods("GET")
	route.Handle("/documentation/{fileUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetFileByFileUUID), false)).Methods("GET")
	route.Handle("/documentation/{fileUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DelFileByFileUUID), true)).Methods("DELETE")

	route.Handle("/generator/start", API.TokenAuthMiddleware(http.HandlerFunc(API.GeneratorStart), true)).Methods("POST")
	route.Handle("/generator/stop", API.TokenAuthMiddleware(http.HandlerFunc(API.GeneratorStop), true)).Methods("POST")
	route.Handle("/generator/equipment/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.GetDataByEquipmentUUID), false)).Methods("GET")
	route.Handle("/generator/equipment/{equipmentUUID}", API.TokenAuthMiddleware(http.HandlerFunc(API.DeleteDataByEquipmentUUID), true)).Methods("DELETE")
	route.Handle("/generator/clear", API.TokenAuthMiddleware(http.HandlerFunc(API.ClearAllData), true)).Methods("DELETE")
	route.HandleFunc("/service/register", API.ServiceRegister).Methods("POST")

	for i := 0; i < 10; i++ {
		fmt.Println(uuid.NewV4())
	}
	if err = addServicesSecretsToDB(Db); err != nil {
		utils.PrintDebug("Error in addServicesSecretsToDB", err.Error())
		return
	}
	address := ":" + conf.ServiceGateway.Port
	utils.PrintDebug("Launched on " + address)

	if err = http.ListenAndServe(address, r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
