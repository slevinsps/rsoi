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
func (db *DataBase) CreateMonitor(monitor models.Monitor) (err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlInsert := `
	INSERT INTO monitor.monitor(name, monitor_uuid, user_uuid) VALUES
    ($1, $2, $3);
		`
	_, err = tx.Exec(sqlInsert, monitor.Name, monitor.MonitorUUID, monitor.UserUUID)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CreateMonitor +")

	return
}

// GetMonitorByMonitorUUIDuserUUID
func (db *DataBase) GetMonitorByMonitorUUIDuserUUID(monitorUID uuid.UUID, userUID uuid.UUID) (monitor models.Monitor, checkFindMonitor bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkFindMonitor = true
	sqlQuery :=
		"SELECT id, name, monitor_uuid, user_uuid FROM monitor.monitor where monitor_uuid = $1 and user_uuid = $2;"

	row := tx.QueryRow(sqlQuery, monitorUID, userUID)
	err = row.Scan(&monitor.ID, &monitor.Name, &monitor.MonitorUUID, &monitor.UserUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindMonitor = false
			err = nil
		}
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetMonitorByUUID +")
	return
}

// GetMonitorByMonitorUUID
func (db *DataBase) GetMonitorByMonitorUUID(monitorUID uuid.UUID) (monitor models.Monitor, checkFindMonitor bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkFindMonitor = true
	sqlQuery :=
		"SELECT id, name, monitor_uuid, user_uuid FROM monitor.monitor where monitor_uuid = $1;"

	row := tx.QueryRow(sqlQuery, monitorUID)
	err = row.Scan(&monitor.ID, &monitor.Name, &monitor.MonitorUUID, &monitor.UserUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindMonitor = false
			err = nil
		}
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetMonitorByUUID +")
	return
}

// GetAllMonitorsByUserUUID
func (db *DataBase) GetAllMonitorsByUserUUID(userUID uuid.UUID) (monitors []models.Monitor, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT id, name, monitor_uuid, user_uuid FROM monitor.monitor where user_uuid = $1;"

	rows, erro := tx.Query(sqlQuery, userUID)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetAllMonitorsByUserUUID Query error")
		return
	}
	for rows.Next() {
		monitor := models.Monitor{}
		if err = rows.Scan(&monitor.ID, &monitor.Name, &monitor.MonitorUUID, &monitor.UserUUID); err != nil {
			utils.PrintDebug("database/GetAllMonitorsByUserUUID wrong row catched")
			break
		}

		monitors = append(monitors, monitor)
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/GetAllMonitorsByUserUUID +")
	return
}

// DelMonitorByMonitorUUID
func (db *DataBase) DelMonitorByMonitorUUID(monitorUID uuid.UUID) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"DELETE FROM monitor.monitor WHERE monitor_uuid = $1;"

	_, err = tx.Exec(sqlQuery, monitorUID)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/DelMonitorByMonitorUUID +")

	return
}

// AddEquipment
func (db *DataBase) AddEquipment(monitorUUID uuid.UUID, equipmentUUID uuid.UUID) (err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlInsert := `
	INSERT INTO monitor.monitor_equipment_xref(monitor_uuid, equipment_uuid) VALUES
    ($1, $2) ON CONFLICT (monitor_uuid, equipment_uuid) DO NOTHING;
		`
	_, err = tx.Exec(sqlInsert, monitorUUID, equipmentUUID)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/AddEquipment +")

	return
}

// DelEquipment
func (db *DataBase) DelEquipment(monitorUUID uuid.UUID, equipmentUUID uuid.UUID) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"DELETE FROM monitor.monitor_equipment_xref WHERE monitor_uuid = $1 and equipment_uuid = $2;"

	_, err = tx.Exec(sqlQuery, monitorUUID, equipmentUUID)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/DelEquipment +")

	return
}
