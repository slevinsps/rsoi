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

// CreateEquipment
func (db *DataBase) CreateEquipment(equipment models.Equipment) (err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlInsert := `
	INSERT INTO equipment.equipment(name, equipment_uuid, equipment_model_uuid, status) VALUES
    ($1, $2, $3, $4);
		`
	_, err = tx.Exec(sqlInsert, equipment.Name, equipment.EquipmentUUID, equipment.EquipmentModelUUID, equipment.Status)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CreateEquipment +")

	return
}

// UpdateEquipment
func (db *DataBase) UpdateEquipment(equipment models.Equipment) (equipment_res models.Equipment, err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery := `
	UPDATE equipment.equipment SET name = $1, equipment_model_uuid = $2, status = $3 where equipment_uuid = $4 RETURNING id, name, equipment_uuid, equipment_model_uuid, status;`
	row := tx.QueryRow(sqlQuery, equipment.Name, equipment.EquipmentModelUUID, equipment.Status, equipment.EquipmentUUID)
	err = row.Scan(&equipment_res.ID, &equipment_res.Name, &equipment_res.EquipmentUUID, &equipment_res.EquipmentModelUUID, &equipment_res.Status)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/UpdateEquipment +")

	return
}

// GetEquipmentByUUID
func (db *DataBase) GetEquipmentByEquipmentUUID(equipmentUID uuid.UUID) (equipment models.Equipment, checkFindEquipment bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkFindEquipment = true
	sqlQuery :=
		"SELECT e.id, e.name, e.equipment_uuid, em.name, e.equipment_model_uuid, e.status FROM equipment.equipment as e JOIN equipment.equipment_model AS em ON em.equipment_model_uuid = e.equipment_model_uuid WHERE e.equipment_uuid = $1;"

	row := tx.QueryRow(sqlQuery, equipmentUID)
	err = row.Scan(&equipment.ID, &equipment.Name, &equipment.EquipmentUUID, &equipment.ModelName, &equipment.EquipmentModelUUID, &equipment.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindEquipment = false
			err = nil
		}
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetEquipmentByUUID +")
	return
}

// GetAllEquipmentsByUserUUID
func (db *DataBase) GetAllEquipments() (equipments []models.Equipment, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT e.id, e.name, em.name, e.equipment_uuid, e.equipment_model_uuid, e.status FROM equipment.equipment as e JOIN equipment.equipment_model AS em ON em.equipment_model_uuid = e.equipment_model_uuid;"

	rows, erro := tx.Query(sqlQuery)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetAllEquipmentsByUserUUID Query error")
		return
	}
	for rows.Next() {
		equipment := models.Equipment{}
		if err = rows.Scan(&equipment.ID, &equipment.Name, &equipment.ModelName, &equipment.EquipmentUUID, &equipment.EquipmentModelUUID, &equipment.Status); err != nil {
			utils.PrintDebug("database/GetAllEquipmentsByUserUUID wrong row catched")
			break
		}

		equipments = append(equipments, equipment)
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/GetAllEquipmentsByUserUUID +")
	return
}

// DelEquipmentByEquipmentUUID
func (db *DataBase) DelEquipmentByEquipmentUUID(equipmentUID uuid.UUID) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"DELETE FROM equipment.equipment WHERE equipment_uuid = $1;"

	_, err = tx.Exec(sqlQuery, equipmentUID)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/DelEquipmentByEquipmentUUID +")

	return
}

// GetEquipmentsByMonitorUUID
func (db *DataBase) GetEquipmentsByMonitorUUID(monitorUUID uuid.UUID) (equipments []models.Equipment, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT e.id, e.name, em.name, e.equipment_uuid, e.equipment_model_uuid, e.status FROM monitor.monitor_equipment_xref as me JOIN equipment.equipment as e ON me.equipment_uuid = e.equipment_uuid JOIN equipment.equipment_model AS em ON em.equipment_model_uuid = e.equipment_model_uuid WHERE me.monitor_uuid = $1;"

	rows, erro := tx.Query(sqlQuery, monitorUUID)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetEquipmentsByMonitorUUID Query error")
		return
	}
	for rows.Next() {
		equipment := models.Equipment{}
		if err = rows.Scan(&equipment.ID, &equipment.Name, &equipment.ModelName, &equipment.EquipmentUUID, &equipment.EquipmentModelUUID, &equipment.Status); err != nil {
			utils.PrintDebug("database/GetEquipmentsByMonitorUUID wrong row catched")
			break
		}

		equipments = append(equipments, equipment)
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/GetEquipmentsByMonitorUUID +")
	return
}

// GetNotAddedEquipmentsByMonitorUUID
func (db *DataBase) GetNotAddedEquipmentsByMonitorUUID(monitorUUID uuid.UUID) (equipments []models.Equipment, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery := "SELECT e.id, e.name, e.equipment_uuid, em.name, e.equipment_model_uuid, e.status FROM equipment.equipment as e JOIN equipment.equipment_model AS em ON em.equipment_model_uuid = e.equipment_model_uuid  WHERE e.id " +
		"not in (SELECT ee.id FROM monitor.monitor_equipment_xref as me join equipment.equipment as ee on me.equipment_uuid = ee.equipment_uuid where me.monitor_uuid = $1);"

	rows, erro := tx.Query(sqlQuery, monitorUUID)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetNotAddedEquipmentsByMonitorUUID Query error")
		return
	}
	for rows.Next() {
		equipment := models.Equipment{}
		if err = rows.Scan(&equipment.ID, &equipment.Name, &equipment.EquipmentUUID, &equipment.ModelName, &equipment.EquipmentModelUUID, &equipment.Status); err != nil {
			utils.PrintDebug("database/GetNotAddedEquipmentsByMonitorUUID wrong row catched")
			break
		}

		equipments = append(equipments, equipment)
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/GetNotAddedEquipmentsByMonitorUUID +")
	return
}
