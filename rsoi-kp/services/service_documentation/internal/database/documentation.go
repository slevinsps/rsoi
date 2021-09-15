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

// CreateFile
func (db *DataBase) CreateFile(file models.File) (err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlInsert := `
	INSERT INTO documentation.documentation(documentation_uuid, equipment_model_uuid, name, path) VALUES
    ($1, $2, $3, $4);
		`
	_, err = tx.Exec(sqlInsert, file.FileUUID, file.EquipmentModelUUID, file.Name, file.Path)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CreateFile +")

	return
}

// UpdateFile
func (db *DataBase) UpdateFile(file models.File) (file_res models.File, err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery := `
	UPDATE documentation.documentation SET path = $1, equipment_model_uuid = $2, name = $3 where documentation_uuid = $4 RETURNING id, documentation_uuid, equipment_model_uuid, name, path;`
	row := tx.QueryRow(sqlQuery, file.Path, file.EquipmentModelUUID, file.Name, file.FileUUID)
	err = row.Scan(&file_res.ID, &file_res.FileUUID, &file_res.EquipmentModelUUID, &file_res.Name, &file_res.Path)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/UpdateFile +")

	return
}

// GetFileByUUID
func (db *DataBase) GetFileByFileUUID(fileUID uuid.UUID) (file models.File, checkFindFile bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkFindFile = true
	sqlQuery :=
		"SELECT id, documentation_uuid, equipment_model_uuid, name, path FROM documentation.documentation where documentation_uuid = $1;"

	row := tx.QueryRow(sqlQuery, fileUID)
	err = row.Scan(&file.ID, &file.FileUUID, &file.EquipmentModelUUID, &file.Name, &file.Path)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindFile = false
			err = nil
		}
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetFileByUUID +")
	return
}

// GetAllFilesByEquipmentModelUUID
func (db *DataBase) GetAllFilesByEquipmentModelUUID(equipmentModelUUID uuid.UUID) (files []models.File, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT id, documentation_uuid, equipment_model_uuid, name, path FROM documentation.documentation where equipment_model_uuid = $1;"

	rows, erro := tx.Query(sqlQuery, equipmentModelUUID)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetAllFilesByEquipmentModelUUID Query error")
		return
	}
	for rows.Next() {
		file := models.File{}
		if err = rows.Scan(&file.ID, &file.FileUUID, &file.EquipmentModelUUID, &file.Name, &file.Path); err != nil {
			utils.PrintDebug("database/GetAllFilesByEquipmentModelUUID wrong row catched")
			break
		}

		files = append(files, file)
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/GetAllFiles +")
	return
}

// GetAllFiles
func (db *DataBase) GetAllFiles() (files []models.File, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT id, documentation_uuid, equipment_model_uuid, name, path FROM documentation.documentation;"

	rows, erro := tx.Query(sqlQuery)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetAllFiles Query error")
		return
	}
	for rows.Next() {
		file := models.File{}
		if err = rows.Scan(&file.ID, &file.FileUUID, &file.EquipmentModelUUID, &file.Name, &file.Path); err != nil {
			utils.PrintDebug("database/GetAllFiles wrong row catched")
			break
		}

		files = append(files, file)
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/GetAllFiles +")
	return
}

// DelFileByFileUUID
func (db *DataBase) DelFileByFileUUID(fileUUID uuid.UUID) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"DELETE FROM documentation.documentation WHERE documentation_uuid = $1;"

	_, err = tx.Exec(sqlQuery, fileUUID)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/DelFileByFileUUID +")

	return
}
