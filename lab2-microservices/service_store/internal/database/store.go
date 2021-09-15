package database

import (
	"database/sql"
	"fmt"
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

// GetUserByUUID
func (db *DataBase) GetUserByUUID(userUID uuid.UUID) (user models.User, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT id, name, user_uid FROM store.users where user_uid = $1;"

	fmt.Println(userUID)
	row := tx.QueryRow(sqlQuery, userUID)
	err = row.Scan(&user.ID, &user.Name, &user.UserUUID)
	if err != nil {
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetUserByUUID +")
	return
}
