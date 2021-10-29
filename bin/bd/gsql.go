package gsql

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	//Драйвер для БД ...

	_ "github.com/mattn/go-sqlite3"
)

//User - объект для работы с польхователями в БД ...
type User struct {
	Login string
	Pass  string
	Email string
}

//NewUser - функция для создания объекта User ...
func NewUser() User {
	return User{}
}

//CreateUser - Создаем запись о пользователе в БД ...
func (us User) CreateUser() (string, error) {
	connect, alert := connection()
	if alert != nil {
		return "", alert
	}

	result, err := connect.Exec("INSERT INTO Users (Name, Pass, Email) VALUES ($1, $2, $3)", us.Login, us.Pass, us.Email)
	if err != nil {
		return "", err
	}

	indx, _ := result.LastInsertId()
	return strconv.Itoa(int(indx)), nil
}

//DeleteUser - Удаляем запись о пользователе в БД ...
func (us User) DeleteUser() error {
	return nil
}

//GetUser - Получаем запись о пользователе в БД ...
func (us User) GetUser() (string, error) {

	uid := ""

	connect, alert := connection()
	if alert != nil {
		return "", alert
	}

	Row := connect.QueryRow("SELECT t.ID from Users as t WHERE t.Name = $1 and t.Pass = $2", us.Login, us.Pass)

	err := Row.Scan(&uid)
	if err != nil {
		return "", err
	}

	log.Println(uid)

	return uid, nil
}

//UpdateUser - Обновляем запись о пользователе в БД ...
func (us User) UpdateUser() error {
	return nil
}

//DBselector - объект для работы с базой данных ...
type DBselector struct {
	ID         string
	User       string
	Login      string
	Pass       string
	Email      string
	Address    string
	CreateDate time.Time
}

func connection() (*sql.DB, error) {
	db, alert := sql.Open("sqlite3", "bin/bd/generDB.db")
	return db, alert
}

//New - создаем экземпляр объекта для работы с бд ...
func New() *DBselector {
	return &DBselector{}
}

//Add - добавление записи в БД ...
func (db DBselector) Add() error {
	connect, alert := connection()
	if alert != nil {
		return alert
	}

	_, err := connect.Exec("INSERT INTO SavedUserData (User, Address, Login, Pass, Email) VALUES ($1, $2, $3, $4, $5)", db.User, db.Address, db.Login, db.Pass, db.Email)
	if err != nil {
		return err
	}

	return nil
}

//Del - удалить запись в БД ...
func (db DBselector) Del() error {
	connect, alert := connection()
	if alert != nil {
		return alert
	}

	_, err := connect.Exec("DELETE from SavedUserData  WHERE User = $1 AND Address = $2", db.User, db.Address)
	if err != nil {
		return err
	}

	return nil
}

//Upd - обновить запись в БД ...
func (db DBselector) Upd() error {
	connect, alert := connection()
	if alert != nil {
		return alert
	}

	_, err := connect.Exec("UPDATE SavedUserData SET Login = $1, Pass = $2, Email= $3 WHERE User = $4 AND Address = $5", db.Login, db.Pass, db.Email, db.User, db.Address)
	if err != nil {
		return err
	}

	return nil
}

//Sel - выборка по параметрам из БД ...
func (db DBselector) Sel() ([]DBselector, error) {

	connect, alert := connection()
	if alert != nil {
		return nil, alert
	}

	Rows, err := connect.Query("SELECT T.User as User, T.Address as Address, T.Login as Login, T.Pass as Pass FROM SavedUserData as T WHERE T.User = $1", db.User)
	if err != nil {
		return nil, err
	}

	selectorArray := []DBselector{}
	userdata := DBselector{}

	for Rows.Next() {
		exeptrow := Rows.Scan(&userdata.User, &userdata.Address, &userdata.Login, &userdata.Pass)
		if exeptrow != nil {
			return nil, exeptrow
		}
		selectorArray = append(selectorArray, userdata)

	}

	return selectorArray, nil
}
