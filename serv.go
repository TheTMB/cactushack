package main

import (
	"net/http"
	"fmt"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"errors"
)

type sign_up_struct struct {
	Login    string
	Password string
}

func handlerSignUp(w http.ResponseWriter, r *http.Request) {
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t sign_up_struct
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for sigh up %v", err)
	}

	user, e := GetUser(t.Login);
	if e != 0 {
		printErr(w, errors.New(fmt.Sprintf("GetUser error %d", e)))
		return
	}
	if (user != nil) {
		u, ok := user.(user_st)
		if !ok {
			printErr(w, errors.New("can't cast to user_st"))
		}
		if ok && u.Login == t.Login {
			printErr(w, errors.New("dublicate login"))
			return
		}
	}

	if err := CreateUser(user_st{Login:t.Login, Passsword:t.Password}); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("create user error %d", err)))
		return;
	}

	w.WriteHeader(http.StatusOK)
}

func handlerLogIn(w http.ResponseWriter, r *http.Request) {
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t sign_up_struct
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for sigh up %v", err)
	}

	user, e := GetUser(t.Login);
	if e != 0 || user == nil {
		printErr(w, errors.New(fmt.Sprintf("GetUser error %d", e)))
		return
	}
	if (user != nil) {
		u, ok := user.(user_st)
		if !ok {
			printErr(w, errors.New("can't cast to user_st"))
			return
		}
		if ok && u.Passsword != t.Password {
			printErr(w, errors.New("wrong password"))
			return
		}
		jsonAnswer := fmt.Sprintf("{\"token\":\"%v\"}", u.ID)
		w.Write([]byte(jsonAnswer))
	}
}

func printErr(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Printf("err with db: %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/signup", handlerSignUp)
	http.HandleFunc("/login", handlerLogIn)
	http.ListenAndServe(":8080", nil)
}