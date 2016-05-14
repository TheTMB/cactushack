package main

import (
	"net/http"
	"fmt"
	"database/sql"
)

var db *sql.DB;

type user_st struct {
	ID       int
	Login    string
	Password string
}

func SetUp() int {
	var err error;
	db, err = sql.Open("mysql", "root:pass@/tmb_db")
	if st, e := checkDbErr(err); st {
		return e;
	}
	return 0
}

func checkDbErr(err error) (status bool, error int) {
	if err != nil {
		fmt.Printf("err with db: %v \n", err)
		status = true;
		error = http.StatusInternalServerError;
	}
	return
}

func GetUser(inputLogin string) (interface{}, int) {
	rows, err := db.Query("SELECT * FROM Users WHERE login=?", inputLogin)
	if st, err := checkDbErr(err); st {
		return nil, err
	}
	if rows.Next() {
		var ID int
		var login string
		var password string
		if st, err := checkDbErr(rows.Scan(&ID, &login, &password)); st {
			return nil, err
		}
		return user_st{ID, login, password}, 0
	}
	return nil, 0
}

func CreateUser(user user_st) (err int) {
	stmt, error := db.Prepare("insert Users SET login=?, password=?")
	if st, err := checkDbErr(error); st {
		return err;
	}

	if _, err := stmt.Exec(user.Login, user.Password); err != nil {
		if st, e := checkDbErr(err); st {
			return e
		}
	}
	return 0
}