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
func (db *DataBase) GetUserByUUID(userUID uuid.UUID) (user models.User, checkFindUser bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkFindUser = true
	sqlQuery :=
		"SELECT id, login, password, user_uuid, is_admin FROM session.users where user_uuid = $1;"

	fmt.Println(userUID)
	row := tx.QueryRow(sqlQuery, userUID)
	err = row.Scan(&user.ID, &user.Login, &user.Password, &user.UserUUID, &user.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindUser = false
			err = nil
		}
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetUserByUUID +")
	return
}

// GetUserByLogin
func (db *DataBase) GetUserByLogin(login string) (user models.User, checkFindUser bool, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkFindUser = true
	sqlQuery :=
		"SELECT id, login, password, user_uuid, is_admin FROM session.users where login = $1;"

	row := tx.QueryRow(sqlQuery, login)
	err = row.Scan(&user.ID, &user.Login, &user.Password, &user.UserUUID, &user.IsAdmin)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindUser = false
			err = nil
		}
		return
	}

	err = tx.Commit()
	utils.PrintDebug("database/GetUserByUUID +")
	return
}

// GetAllUsers
func (db *DataBase) GetAllUsers() (users []models.User, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT id, login, password, user_uuid, is_admin  FROM session.users;"

	rows, erro := tx.Query(sqlQuery)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/isNicknameUnique_test Query error")
		return
	}
	for rows.Next() {
		user := models.User{}
		if err = rows.Scan(&user.ID, &user.Login, &user.Password, &user.UserUUID, &user.IsAdmin); err != nil {
			utils.PrintDebug("database/isNicknameUnique_test wrong row catched")
			break
		}

		users = append(users, user)
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/GetAllUsers +")
	return
}

// CreateUser
func (db *DataBase) CreateUser(user models.User) (checkFindUser bool, err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	checkFindUser = false
	if _, checkFindUser, err = db.GetUserByLogin(user.Login); err != nil {
		utils.PrintDebug("database/CreateUser - fail uniqie")
		return
	}

	if checkFindUser {
		utils.PrintDebug("CreateUser ", user)
		return
	}

	//utils.PrintDebug(user)
	sqlInsert := `
	INSERT INTO session.users(login, password, user_uuid, is_admin) VALUES
    ($1, $2, $3, $4) on conflict do nothing;
		`
	_, err = tx.Exec(sqlInsert, user.Login, user.Password, user.UserUUID, user.IsAdmin)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/CreateUser +")

	return
}
