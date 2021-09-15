package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"services/config"
	"services/utils"

	"github.com/go-redis/redis/v7"
	_ "github.com/lib/pq"
)

func DataBaseInitialize(confPath string) (db *sql.DB, redisClient *redis.Client, conf *config.Configuration, err error) {

	if conf, err = config.Init(confPath); err != nil {
		return
	}

	if db, err = Init(conf.DataBase); err != nil {
		return
	}

	if redisClient, err = InitRedis(conf.Redis); err != nil {
		return
	}

	return
}

func InitRedis(redisConfig config.RedisConfig) (client *redis.Client, err error) {
	address := redisConfig.Host + ":" + redisConfig.Port
	fmt.Println(address)
	client = redis.NewClient(&redis.Options{
		Addr:       address,
		Password:   "",
		DB:         0,
		MaxRetries: 3,
	})

	_, err = client.Ping().Result()
	return
}

func Init(CDB config.DatabaseConfig) (database *sql.DB, err error) {

	fmt.Println("DATABASE_URL = ", os.Getenv(CDB.URL))
	// for local launch
	if os.Getenv(CDB.URL) == "" {
		os.Setenv(CDB.URL, CDB.LocalDataBaseUrl)
	}

	fmt.Println("DATABASE_URL = ", os.Getenv(CDB.URL))

	if database, err = sql.Open(CDB.DriverName, "postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"); err != nil {
		utils.PrintDebug("database/Init cant open:" + err.Error())
		return
	}

	database.SetMaxOpenConns(CDB.MaxOpenConns)

	if err = database.Ping(); err != nil {
		utils.PrintDebug("database/Init cant access: " + err.Error())
		return
	}

	utils.PrintDebug("database/Init open")
	return
}

func (db *DataBase) CreateTables() error {
	query, err := ioutil.ReadFile("init.pgsql")
	if err != nil {
		panic(err)
	}
	utils.PrintDebug("here1")
	_, err = db.Db.Exec(string(query))
	utils.PrintDebug("here2")

	if err != nil {
		utils.PrintDebug("database/init - fail:" + err.Error())
	}
	return err
}
