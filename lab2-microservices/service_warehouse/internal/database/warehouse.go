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

// startWarrantyPeriod
func (db *DataBase) InsertOrderItem(orderItemTableRecord models.OrderItemTableRecord) (err error) {
	var (
		tx *sql.Tx
	)

	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlString := `
	INSERT INTO warehouse.order_item(canceled, order_item_uid, order_uid, item_id) VALUES
    ($1, $2, $3, $4);
		`

	_, err = tx.Exec(sqlString, orderItemTableRecord.Canceled, orderItemTableRecord.OrderItemUid, orderItemTableRecord.OrderUid, orderItemTableRecord.ItemID)
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

// GetItemInfoByModeSize
func (db *DataBase) GetItemInfoByModeSize(model string, size string) (item models.Item, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT id, model, size, available_count FROM warehouse.items WHERE model = $1 and size = $2;"

	row := tx.QueryRow(sqlQuery, model, size)
	err = row.Scan(&item.ID, &item.Model, &item.Size, &item.AvailableCount)

	if err != nil {
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetItemInfoByModeSize +")
	return
}

// GetItemInfoByOrderItemUID
func (db *DataBase) GetItemInfoByOrderItemUID(orderItemUID uuid.UUID, availableCount bool) (item models.Item, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	count := 0
	sqlQuery :=
		"SELECT items.model, items.size, items.available_count FROM warehouse.items items JOIN warehouse.order_item  order_item ON order_item.item_id = items.id Where order_item.order_item_uid = $1;"

	row := tx.QueryRow(sqlQuery, orderItemUID)
	if availableCount {
		err = row.Scan(&item.Model, &item.Size, &item.AvailableCount)
	} else {
		err = row.Scan(&item.Model, &item.Size, &count)
	}

	if err != nil {
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetItemInfoByOrderItemUID +")
	return
}

// UpdateItem
func (db *DataBase) TakeOneItem(item models.Item) (err error) {
	var (
		tx *sql.Tx
	)

	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlString := `
	UPDATE warehouse.items SET available_count = available_count - 1 WHERE id = $1;
		`

	_, err = tx.Exec(sqlString, item.ID)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/UpdateItem +")

	return
}

// ReturnOneItem
func (db *DataBase) ReturnOneItem(orderItemUID uuid.UUID) (err error) {
	var (
		tx *sql.Tx
	)

	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlString := `
	UPDATE warehouse.items SET available_count = available_count + 1 ` +
		`where id = (select item_id from warehouse.order_item where order_item_uid = $1)`

	_, err = tx.Exec(sqlString, orderItemUID)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/ReturnOneItem +")

	return
}

// CancelOrderItem
func (db *DataBase) CancelOrderItem(orderItemUID uuid.UUID) (err error) {
	var (
		tx *sql.Tx
	)

	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlString := `
	UPDATE warehouse.order_item SET canceled = true where order_item_uid = $1`
	_, err = tx.Exec(sqlString, orderItemUID)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CancelOrderItem +")

	return
}
