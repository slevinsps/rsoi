package config

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	ServiceSession       ServerConfig   `json:"service_session"`
	ServiceMonitor       ServerConfig   `json:"service_monitor"`
	ServiceEquipment     ServerConfig   `json:"service_equipment"`
	ServiceDocumentation ServerConfig   `json:"service_documentation"`
	ServiceGenerator     ServerConfig   `json:"service_generator"`
	ServiceGateway       ServerConfig   `json:"service_gateway"`
	DataBase             DatabaseConfig `json:"dataBase"`
	Redis                RedisConfig    `json:"redis_gateway"`
}

type RedisConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	URL  string `json:"url"`
}

type DatabaseConfig struct {
	DriverName       string `json:"driverName"`
	URL              string `json:"url"`
	MaxOpenConns     int    `json:"maxOpenConns"`
	LocalDataBaseUrl string `json:"localDataBaseUrl"`
}

type DlocalDatabaseConfig struct {
	DriverName   string `json:"driverName"`
	URL          string `json:"url"`
	MaxOpenConns int    `json:"maxOpenConns"`
}

// Init load configuration file
func Init(path string) (conf *Configuration, err error) {
	conf = &Configuration{}
	var data []byte

	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	err = json.Unmarshal(data, conf)

	return
}
