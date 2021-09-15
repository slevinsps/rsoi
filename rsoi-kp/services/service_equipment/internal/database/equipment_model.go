package database

import (
	"database/sql"
	"services/internal/models"
	"services/utils"

	uuid "github.com/satori/go.uuid"
)

// CreateEquipmentModel
func (db *DataBase) CreateEquipmentModel(equipmentModel models.EquipmentModel) (err error) {

	var (
		tx *sql.Tx
	)
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlInsert := `
	INSERT INTO equipment.equipment_model(equipment_model_uuid, name) VALUES
    ($1, $2);
		`

	_, err = tx.Exec(sqlInsert, equipmentModel.EquipmentModelUUID, equipmentModel.Name)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CreateEquipmentModel +")

	return
}

// GetEquipmentModelByUUID
func (db *DataBase) GetEquipmentModelByEquipmentModelUUID(equipmentModelUID uuid.UUID) (equipmentModel models.EquipmentModel, checkFindEquipmentModel bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkFindEquipmentModel = true
	sqlQuery :=
		"SELECT id, equipment_model_uuid, name FROM equipment.equipment_model where equipment_model_uuid = $1;"

	row := tx.QueryRow(sqlQuery, equipmentModelUID)
	err = row.Scan(&equipmentModel.ID, &equipmentModel.EquipmentModelUUID, &equipmentModel.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindEquipmentModel = false
			err = nil
		}
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetEquipmentModelByUUID +")
	return
}

// GetAllEquipmentModelsByUserUUID
func (db *DataBase) GetAllEquipmentModels() (equipmentModels []models.EquipmentModel, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT id, equipment_model_uuid, name FROM equipment.equipment_model;"

	rows, erro := tx.Query(sqlQuery)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetAllEquipmentModelsByUserUUID Query error")
		return
	}
	for rows.Next() {
		equipmentModel := models.EquipmentModel{}
		if err = rows.Scan(&equipmentModel.ID, &equipmentModel.EquipmentModelUUID, &equipmentModel.Name); err != nil {
			utils.PrintDebug("database/GetAllEquipmentModelsByUserUUID wrong row catched")
			break
		}

		equipmentModels = append(equipmentModels, equipmentModel)
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/GetAllEquipmentModelsByUserUUID +")
	return
}

// DelEquipmentModelByEquipmentModelUUID
func (db *DataBase) DelEquipmentModelByEquipmentModelUUID(equipmentModelUID uuid.UUID) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"DELETE FROM equipment.equipment_model WHERE equipment_model_uuid = $1;"

	_, err = tx.Exec(sqlQuery, equipmentModelUID)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/DelEquipmentModelByEquipmentModelUUID +")

	return
}
