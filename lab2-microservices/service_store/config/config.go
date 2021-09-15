package config

import (
	"encoding/json"
	"io/ioutil"
)

type Configuration struct {
	ServiceOrders    ServerConfig   `json:"service_orders"`
	ServiceWarranty  ServerConfig   `json:"service_warranty"`
	ServiceStore     ServerConfig   `json:"service_store"`
	ServiceWarehouse ServerConfig   `json:"service_warehouse"`
	DataBase         DatabaseConfig `json:"dataBase"`
}

type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	URL  string `json:"url"`
}

// DatabaseConfig set type of database management system
//   the url of connection string, max amount of
//   connections, tables, sizes of page  of gamers
//   and users
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
