package database

import (
	"database/sql"
	"fmt"
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

// OrderCreate
func (db *DataBase) OrderCreate(orderUID uuid.UUID, userUID uuid.UUID, itemUID uuid.UUID) (err error) {
	var (
		tx         *sql.Tx
		id         int
		item_uid   uuid.UUID
		order_date time.Time
		order_uid  uuid.UUID
		status_    string
		user_uid   uuid.UUID
	)

	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlInsert := `
	INSERT INTO orders.orders(item_uid, order_date, order_uid, status, user_uid) VALUES
    ($1, $2, $3, $4, $5) RETURNING id, item_uid, order_date, order_uid, status, user_uid;
		`
	currentTime := time.Now()
	date := currentTime.Format("2006-01-02")
	status := "PAID"
	row := tx.QueryRow(sqlInsert, itemUID, date, orderUID, status, userUID)

	err = row.Scan(&id, &item_uid, &order_date, &order_uid, &status_, &user_uid)
	if err != nil {
		return
	}
	fmt.Println(id)
	fmt.Println(item_uid)
	fmt.Println(order_date)
	fmt.Println(order_uid)
	fmt.Println(status_)
	fmt.Println(user_uid)

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/PersonCreate +")

	return
}

// GetUserOrderInfo
func (db *DataBase) GetUserOrderInfo(userUID uuid.UUID, orderUID uuid.UUID) (order models.Order, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT order_uid, order_date, item_uid, status FROM orders.orders where user_uid = $1 and order_uid = $2;"

	row := tx.QueryRow(sqlQuery, userUID, orderUID)
	err = row.Scan(&order.OrderUID, &order.OrderDate, &order.ItemUID, &order.Status)

	if err != nil {
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetUserOrderInfo +")
	return
}

// GetOrderInfoByOrderUID
func (db *DataBase) GetOrderInfoByOrderUID(orderUID uuid.UUID) (order models.Order, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT order_uid, order_date, item_uid, status FROM orders.orders where order_uid = $1;"

	row := tx.QueryRow(sqlQuery, orderUID)
	err = row.Scan(&order.OrderUID, &order.OrderDate, &order.ItemUID, &order.Status)

	if err != nil {
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetOrderInfoByOrderUID +")
	return
}

// GetUserOrdersInfo
func (db *DataBase) GetUserOrdersInfo(userUID uuid.UUID) (orders []models.Order, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT order_uid, order_date, item_uid, status FROM orders.orders where user_uid = $1;"

	rows, erro := tx.Query(sqlQuery, userUID)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetUserOrdersInfo Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		order := models.Order{}
		if err = rows.Scan(&order.OrderUID, &order.OrderDate, &order.ItemUID, &order.Status); err != nil {
			utils.PrintDebug("database/GetUserOrdersInfo wrong row catched")
			break
		}
		orders = append(orders, order)
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.PrintDebug("database/GetUserOrdersInfo +")

	return
}
