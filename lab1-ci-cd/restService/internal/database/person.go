package database

import (
	"database/sql"
	"restService/internal/models"
	"restService/internal/utils"
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

// PersonCreate
func (db *DataBase) PersonCreate(person models.Person) (createdPerson models.Person, err error) {

	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlInsert := `
	INSERT INTO Persons(name, age, address, work) VALUES
    ($1, $2, $3, $4) RETURNING id, name, age, address, work;
		`

	row := tx.QueryRow(sqlInsert, person.Name, person.Age, person.Address, person.Work)
	err = row.Scan(&createdPerson.ID, &createdPerson.Name, &createdPerson.Age, &createdPerson.Address, &createdPerson.Work)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/PersonCreate +")

	return
}

// GetAllPersonsInfo
func (db *DataBase) GetAllPersonsInfo() (persons []models.Person, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT id, name, age, address, work FROM Persons ORDER BY id;"

	rows, erro := tx.Query(sqlQuery)
	if erro != nil {
		err = erro
		utils.PrintDebug("database/GetAllPersonsInfo Query error")
		return
	}

	defer rows.Close()

	for rows.Next() {
		person := models.Person{}
		if err = rows.Scan(&person.ID, &person.Name, &person.Age, &person.Address, &person.Work); err != nil {
			utils.PrintDebug("database/GetAllPersonsInfo wrong row catched")
			break
		}
		persons = append(persons, person)
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	utils.PrintDebug("database/GetAllPersonsInfo +")

	return
}

// GetPersonByID
func (db *DataBase) GetPersonByID(id int) (person models.Person, checkFindPerson bool, err error) {
	checkFindPerson = true
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"SELECT id, name, age, address, work FROM Persons " +
			"where id = $1;"

	row := tx.QueryRow(sqlQuery, id)
	err = row.Scan(&person.ID, &person.Name, &person.Age, &person.Address, &person.Work)
	if err != nil {
		if err == sql.ErrNoRows {
			checkFindPerson = false
			err = nil
		}
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}

	utils.PrintDebug("database/GetPersonByID +")

	return
}

// UpdatePersonInfo
func (db *DataBase) UpdatePersonInfo(person models.Person) (updatePerson models.Person, err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery := `
		UPDATE Persons SET name = $1, age = $2, address = $3, work = $4 WHERE id = $5
		`
	_, err = tx.Exec(sqlQuery, person.Name, person.Age, person.Address, person.Work, person.ID)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	utils.PrintDebug("database/UpdatePersonInfo +")

	return
}

// DeletePersonInfo
func (db *DataBase) DeletePersonInfo(id int) (err error) {
	var tx *sql.Tx
	tx, err = db.Db.Begin()
	defer tx.Rollback()

	sqlQuery :=
		"DELETE from Persons where id = $1;"

	_, err = tx.Exec(sqlQuery, id)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}

	utils.PrintDebug("database/DeletePersonInfo +")

	return
}
