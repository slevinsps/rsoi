package database

import (
	"database/sql"
	"services/internal/models"
	"services/utils"
)

type DataBase struct {
	Db *sql.DB
}

func NewDataBase(db *sql.DB) (database *DataBase, err error) {

	database = &DataBase{
		Db: db,
	}

	err = database.CreateTables()
	return
}

// GetServiceByLogin
func (db *DataBase) GetServiceByLogin(login string) (service models.Service, checkFindService bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkFindService = true
	sqlQuery :=
		"SELECT login, password FROM gateway.services_secrets where login = $1;"

	row := tx.QueryRow(sqlQuery, login)
	err = row.Scan(&service.Login, &service.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindService = false
			err = nil
		}
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetServiceByLogin +")
	return
}

// AddServiceSecrets
func (db *DataBase) AddServiceSecrets(service models.Service) (err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	//utils.PrintDebug(user)
	sqlInsert := `
	INSERT INTO gateway.services_secrets(login, password) VALUES
    ($1, $2) on conflict do nothing;
		`
	_, err = tx.Exec(sqlInsert, service.Login, service.Password)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/AddServiceSecrets +")

	return
}
