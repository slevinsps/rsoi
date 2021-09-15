package main

import (
	"database/sql"
	"os"
	"services/config"
	"services/internal/database"
	"services/internal/models"
	"services/internal/repInterface"
	api "services/internal/services"
	utils "services/utils"

	"net/http"

	"github.com/go-redis/redis/v7"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

func main() {
	const confPath = "conf.json"

	var (
		db          *sql.DB
		redisClient *redis.Client
		conf        *config.Configuration
		err         error
		admin       models.User
	)

	if err = godotenv.Load(); err != nil {
		utils.PrintDebug("Error loading .env file")
		return
	}

	if db, redisClient, conf, err = database.DataBaseInitialize(confPath); err != nil {
		utils.PrintDebug("Error in configuration file or database " + err.Error())
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
	route := r.PathPrefix("/api/v1/session").Subrouter()
	route.HandleFunc("/signup", API.CreateUser).Methods("POST")
	route.HandleFunc("/signin", API.SignIn).Methods("POST")
	route.Handle("/logout", API.TokenAuthMiddleware(http.HandlerFunc(API.Logout))).Methods("POST")
	route.HandleFunc("/token/refresh", API.RefreshTokens).Methods("POST")

	route.HandleFunc("/users", API.GetAllUsers).Methods("GET")
	route.HandleFunc("/user/{userUid}", API.GetUserByUUID).Methods("GET")
	route.HandleFunc("/check/user", API.GetUserByToken).Methods("GET")
	route.HandleFunc("/service/register", API.ServiceRegister).Methods("POST")

	// address := conf.ServiceMonitor.Host + ":" + conf.ServiceMonitor.Port
	address := ":" + conf.ServiceSession.Port
	utils.PrintDebug("Launched on " + address)

	admin.Login = os.Getenv("ADMIN_LOGIN")
	admin.UserUUID = uuid.NewV4()
	if admin.Password, err = api.HashPassword(os.Getenv("ADMIN_PASSWORD")); err != nil {
		utils.PrintDebug("Error in HashPassword admin " + err.Error())
	}
	admin.IsAdmin = true
	if _, err = Db.CreateUser(admin); err != nil {
		utils.PrintDebug("Error in create admin " + err.Error())
	}

	if err = http.ListenAndServe(address, r); err != nil {
		utils.PrintDebug("Error:" + err.Error())
	}
}
