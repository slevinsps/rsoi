package database

import (
	"database/sql"
	"fmt"
	"services/constants"
	"services/internal/models"
	"services/utils"
	"time"

	uuid "github.com/satori/go.uuid"
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

// startWarrantyPeriod
func (db *DataBase) StartWarrantyPeriod(itemUID uuid.UUID) (err error) {
	var (
		tx            *sql.Tx
		id            int
		item_uid      uuid.UUID
		warranty_date time.Time
		status_       string
	)

	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlString := `
	INSERT INTO warranty.warranty(item_uid, status, warranty_date) VALUES
    ($1, $2, $3) RETURNING id, item_uid, status, warranty_date;
		`
	currentTime := time.Now()
	date := currentTime
	status := constants.WARRANTY_STATUS_ON
	row := tx.QueryRow(sqlString, itemUID, status, date)

	err = row.Scan(&id, &item_uid, &status_, &warranty_date)
	if err != nil {
		return
	}
	fmt.Println(id)
	fmt.Println(item_uid)
	fmt.Println(status_)
	fmt.Println(warranty_date)

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/startWarrantyPeriod +")

	return
}

// GetUserOrderInfo
func (db *DataBase) GetWarranty(itemUID uuid.UUID) (warranty models.Warranty, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT item_uid, warranty_date, status FROM warranty.warranty where item_uid = $1;"

	row := tx.QueryRow(sqlQuery, itemUID)
	err = row.Scan(&warranty.ItemUID, &warranty.WarrantyDate, &warranty.Status)

	if err != nil {
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetUserOrderInfo +")
	return
}

// UpdateWarranty
func (db *DataBase) UpdateWarranty(warranty models.Warranty) (err error) {
	var (
		tx *sql.Tx
	)

	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlString := `
	UPDATE warranty.warranty SET (comment_, status) = ($2, $3) WHERE item_uid = $1;
		`

	_, err = tx.Exec(sqlString, warranty.ItemUID, warranty.Comment, warranty.Status)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/startWarrantyPeriod +")

	return
}

// CloseWarranty
func (db *DataBase) CloseWarranty(itemUID uuid.UUID) (err error) {
	var (
		tx *sql.Tx
	)

	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlString := `
	DELETE FROM warranty.warranty WHERE item_uid = $1;
		`

	_, err = tx.Exec(sqlString, itemUID)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CloseWarranty +")

	return
}
