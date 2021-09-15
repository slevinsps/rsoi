package database

import (
	"database/sql"
	"services/internal/models"
	"services/utils"

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

// CreateMonitor
func (db *DataBase) CreateData(data models.Data) (err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlInsert := `
	INSERT INTO generator.generator(data_uuid, equipment_uuid, temperature, voltage, frequency, load_level, timestamp_) VALUES
    ($1, $2, $3, $4, $5, $6, $7);
		`
	_, err = tx.Exec(sqlInsert, data.DataUUID, data.EquipmentUUID, data.Temperature, data.Voltage, data.Frequency, data.LoadLevel, data.Timestamp)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CreateData +")

	return
}

// GetAllDatasByUserUUID
func (db *DataBase) GetDataByEquipmentUUID(equipmentUUID uuid.UUID) (data_arr []models.Data, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT id, data_uuid, equipment_uuid, temperature, voltage, frequency, load_level, timestamp_ FROM generator.generator where equipment_uuid = $1 ORDER BY timestamp_ DESC LIMIT 1000;"

	rows, erro := tx.Query(sqlQuery, equipmentUUID)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetAllDatasByUserUUID Query error")
		return
	}
	for rows.Next() {
		data := models.Data{}
		if err = rows.Scan(&data.ID, &data.DataUUID, &data.EquipmentUUID, &data.Temperature, &data.Voltage, &data.Frequency, &data.LoadLevel, &data.Timestamp); err != nil {
			utils.PrintDebug("database/GetAllDatasByUserUUID wrong row catched")
			break
		}

		data_arr = append(data_arr, data)
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/GetAllDatasByUserUUID +")
	return
}

// DeleteDataByEquipmentUUID
func (db *DataBase) DeleteDataByEquipmentUUID(equipmentUUID uuid.UUID) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"DELETE FROM generator.generator WHERE equipment_uuid = $1;"

	_, err = tx.Exec(sqlQuery, equipmentUUID)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/DeleteDataByEquipmentUUID +")

	return
}

// CLearData
func (db *DataBase) CLearData() (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"TRUNCATE generator.generator;"

	_, err = tx.Exec(sqlQuery)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CLearData +")

	return
}
