package main

import (
	"net/http"
	"fmt"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"errors"
)

type user_sign_up_struct struct {
	Login    string
	Password string
}

type user_add_info_struct struct {
	Token            int
	FirstName        string
	LastName         string
	Country          string
	City             string
	University       string
	Start_study      int64
	End_study        int64
	Age              int
	Work             string
	Known_technology string
	About            string
}

type user_get_info_struct struct {
	Token int
}

type company_sign_up_struct struct {
	ID                  int
	Login               string
	Password            string
	FirstName           string
	LastName            string
	Country             string
	City                string
	Phone               string
	Site_domain_address string
	Description         string
}

func handlerSignUpUser(w http.ResponseWriter, r *http.Request) {
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t user_sign_up_struct
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

	if err := CreateUser(user_st{Login:t.Login, Password:t.Password}); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("create user error %d", err)))
		return;
	}

	w.WriteHeader(http.StatusOK)
}

func handlerSignInUser(w http.ResponseWriter, r *http.Request) {
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t user_sign_up_struct
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for sigh up %v", err)
		printErr(w, err)
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
		if ok && u.Password != t.Password {
			printErr(w, errors.New("wrong password"))
			return
		}
		jsonAnswer := fmt.Sprintf("{\"token\":\"%v\"}", u.ID)
		w.Write([]byte(jsonAnswer))
	}
}

func handlerAddUserInfo(w http.ResponseWriter, r *http.Request) {
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t user_add_info_struct
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for sigh up %v", err)
		printErr(w, err)
	}

	if (t.Token == 0) {
		printErr(w, errors.New("token is empty"))
	}

	user, e := GetUserByToken(t.Token);
	if e != 0 || user == nil {
		printErr(w, errors.New(fmt.Sprintf("GetUser error %d", e)))
		return
	}

	if err := UpdateUser(user_st{t.Token, "", "", t.FirstName, t.LastName, t.Country, t.City, t.University, t.Start_study, t.End_study,
		t.Age, t.Work, t.Known_technology, t.About}); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("create user error %d", err)))
		return;
	}

	w.WriteHeader(http.StatusOK)
}

func handlerGetInfoUser(w http.ResponseWriter, r *http.Request) {
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t user_get_info_struct
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for sigh up %v", err)
		printErr(w, err)
	}

	if (t.Token == 0) {
		printErr(w, errors.New("token is empty"))
	}
	fmt.Print("token: " + t.Token)
	user, e := GetUserByToken(t.Token);
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

		jsonAnswer := fmt.Sprintf("{\"token\":\"%v\"}", u.ID)
		w.Write([]byte(jsonAnswer))
	}

}

func handlerSignUpCompany(w http.ResponseWriter, r *http.Request) {
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t company_sign_up_struct
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for sigh up %v", err)
	}

	user, e := GetCompany(t.Login);
	if e != 0 {
		printErr(w, errors.New(fmt.Sprintf("GetUser error %d", e)))
		return
	}
	if (user != nil) {
		c, ok := user.(company_st)
		if !ok {
			printErr(w, errors.New("can't cast to company_st"))
		}
		if ok && c.Login == t.Login {
			printErr(w, errors.New("dublicate login"))
			return
		}
	}

	if err := CreateCompany(company_st{Login:t.Login, Password:t.Password, FirstName:t.FirstName, LastName:t.LastName,
		Country:t.Country, City:t.City, Phone:t.Phone, Site_domain_address:t.Site_domain_address, Description:t.Description}); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("create company error %d", err)))
		return;
	}

	w.WriteHeader(http.StatusOK)
}

func printErr(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Printf("err with db: %v \n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/user_signup", handlerSignUpUser)
	http.HandleFunc("/user_signin", handlerSignInUser)
	http.HandleFunc("/user_addInfo", handlerAddUserInfo)
	http.HandleFunc("/user_getInfo", handlerGetInfoUser)
	http.HandleFunc("/company_signup", handlerSignUpCompany)
	http.ListenAndServe(":8080", nil)
}