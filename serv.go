package main

import (
	"net/http"
	"fmt"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"errors"
	"strconv"
)

type user_sign_struct struct {
	Login    string
	Password string
}

type company_sign_struct struct {
	Login    string
	Password string
}

type user_add_info_struct struct {
	Token            string
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

type user_access struct {
	Token string
}

type company_access struct {
	Token string
}

type company_get_events struct {
	Company_id int
}

type user_get_events struct {
	User_id int
}

type user_join_event_st struct {
	Token    string
	Event_id int
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

type company_add_event_struct struct {
	Name        string
	Description string
	Start_event int64
	End_event   int64
	Company_id  string
	Token       string
}

func handlerSignUpUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t user_sign_struct
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
}

func handlerSignInUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t user_sign_struct
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
		jsonAnswer := fmt.Sprintf("{\"token\":\"%v_u\", \"login\":\"%v\"}", u.ID, u.Login)
		w.Write([]byte(jsonAnswer))
	}
}

func handlerAddUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

	if t.Token == "" || len(t.Token) < 3 {
		printErr(w, errors.New("token less than 3 symbols"))
		return
	}
	token, err := strconv.Atoi(t.Token[:len(t.Token) - 2])
	if err != nil {
		printErr(w, errors.New("token is not int"))
		return
	}
	if (token == 0) {
		printErr(w, errors.New("token is empty"))
		return
	}

	user, e := GetUserByToken(token);
	if e != 0 || user == nil {
		printErr(w, errors.New(fmt.Sprintf("GetUser error %d", e)))
		return
	}

	if err := UpdateUser(user_st{token, "", "", t.FirstName, t.LastName, t.Country, t.City, t.University, t.Start_study, t.End_study,
		t.Age, t.Work, t.Known_technology, t.About}); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("create user error %d", err)))
		return;
	}

	w.WriteHeader(http.StatusOK)
}

func handlerGetInfoUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t user_access
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for sigh up %v", err)
		printErr(w, err)
		return
	}

	if t.Token == "" || len(t.Token) < 3 {
		printErr(w, errors.New("token less than 3 symbols"))
		return
	}
	token, err := strconv.Atoi(t.Token[:len(t.Token) - 2])
	if err != nil {
		printErr(w, errors.New("token is not int"))
	}
	user, e := GetUserByToken(token);
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

		jsonAnswer, err := json.Marshal(u)
		if err != nil {
			printErr(w, err)
			return
		}
		w.Write([]byte(jsonAnswer))
	}

}

func handlerSignUpCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

func handlerSignInCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t company_sign_struct
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for sigh up %v", err)
		printErr(w, err)
	}

	user, e := GetCompany(t.Login);
	if e != 0 || user == nil {
		printErr(w, errors.New(fmt.Sprintf("GetCompany error %d", e)))
		return
	}
	if (user != nil) {
		c, ok := user.(company_st)
		if !ok {
			printErr(w, errors.New("can't cast to company_st"))
			return
		}
		if ok && c.Password != t.Password {
			printErr(w, errors.New("wrong password"))
			return
		}
		jsonAnswer := fmt.Sprintf("{\"token\":\"%v_c\"}", c.ID)
		w.Write([]byte(jsonAnswer))
	}
}

func handlerGetInfoCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t company_access
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for get compnay user %v", err)
		printErr(w, err)
		return
	}
	if t.Token == "" || len(t.Token) < 3 {
		printErr(w, errors.New("token less than 3 symbols"))
		return
	}
	token, err := strconv.Atoi(t.Token[:len(t.Token) - 2])
	if err != nil {
		printErr(w, errors.New("token is not int"))
		return
	}
	if (token == 0) {
		printErr(w, errors.New("token is empty"))
		return
	}
	company, e := GetCompanyByToken(token);
	if e != 0 || company == nil {
		printErr(w, errors.New(fmt.Sprintf("GetCompany error %d", e)))
		return
	}
	if (company != nil) {
		c, ok := company.(company_st)
		if !ok {
			printErr(w, errors.New("can't cast to company_st"))
			return
		}

		jsonAnswer, err := json.Marshal(c)
		if err != nil {
			printErr(w, err)
			return
		}
		w.Write([]byte(jsonAnswer))
	}
}

func handlerAddEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t company_add_event_struct
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for get compnay user %v", err)
		printErr(w, err)
		return
	}
	if t.Token == "" || len(t.Token) < 3 {
		printErr(w, errors.New("token less than 3 symbols"))
		return
	}
	tokenCode := t.Token[len(t.Token) - 2:]
	if tokenCode != "_c" {
		printErr(w, errors.New("this is not company code"))
		return
	}
	token, err := strconv.Atoi(t.Token[:len(t.Token) - 2])
	if err != nil {
		printErr(w, errors.New("token is not int"))
		return
	}
	if (token == 0) {
		printErr(w, errors.New("token is empty"))
		return
	}

	if err := CreateEvent(company_event_st{Name:t.Name, Description:t.Description, Start_event:t.Start_event,
		End_event:t.End_event, Token:token}); err != 0 {
		printErr(w, errors.New("create event error %d"))
	}
}

func handleGetAllEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t user_access
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for get compnay user %v", err)
		printErr(w, err)
		return
	}
	if t.Token == "" || len(t.Token) < 3 {
		printErr(w, errors.New("token less than 3 symbols"))
		return
	}
	token, err := strconv.Atoi(t.Token[:len(t.Token) - 2])
	if err != nil {
		printErr(w, errors.New("token is not int"))
		return
	}
	if (token == 0) {
		printErr(w, errors.New("token is empty"))
		return
	}

	events, e := GetEvents();
	if e != 0 || events == nil {
		printErr(w, errors.New(fmt.Sprintf("GetAllEvents error %d", e)))
		return
	}

	if (events != nil) {
		jsonAnswer, err := json.Marshal(events)
		if err != nil {
			printErr(w, err)
			return
		}
		w.Write([]byte(jsonAnswer))
	}
}

func handleGetCompanyEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t company_get_events
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for get compnay user %v", err)
		printErr(w, err)
		return
	}

	if (t.Company_id == 0) {
		printErr(w, errors.New("token is empty"))
		return
	}

	events, e := GetEventsByCompany(t.Company_id);
	if e != 0 || events == nil {
		printErr(w, errors.New(fmt.Sprintf("GetAllEvents error %d", e)))
		return
	}

	if (events != nil) {
		jsonAnswer, err := json.Marshal(events)
		if err != nil {
			printErr(w, err)
			return
		}
		w.Write([]byte(jsonAnswer))
	}
}

func handleUserJoinEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t user_join_event_st
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for sigh up %v", err)
		printErr(w, err)
	}
	if t.Token == "" || len(t.Token) < 3 {
		printErr(w, errors.New("token less than 3 symbols"))
		return
	}
	token, err := strconv.Atoi(t.Token[:len(t.Token) - 2])
	if err != nil {
		printErr(w, errors.New("token is not int"))
		return
	}
	if (token == 0) {
		printErr(w, errors.New("token is empty"))
		return
	}

	if err := CreateUserEvent(company_user_event_st{user_id:token, event_id:t.Event_id}); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("create user error %d", err)))
		return;
	}

	w.WriteHeader(http.StatusOK)
}

func handleGetUserEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := SetUp(); err != 0 {
		printErr(w, errors.New(fmt.Sprintf("cannot setUp db %d", err)))
	}
	decoder := json.NewDecoder(r.Body)
	var t user_get_events
	err := decoder.Decode(&t)
	if err != nil {
		fmt.Printf("error parse json for get compnay user %v", err)
		printErr(w, err)
		return
	}

	if (t.User_id == 0) {
		printErr(w, errors.New("token is empty"))
		return
	}

	events, e := GetEventsByUser(t.User_id);
	if e != 0 || events == nil {
		printErr(w, errors.New(fmt.Sprintf("GetAllEvents error %d", e)))
		return
	}

	if (events != nil) {
		jsonAnswer, err := json.Marshal(events)
		if err != nil {
			printErr(w, err)
			return
		}
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
	http.HandleFunc("/user_signup", handlerSignUpUser)
	http.HandleFunc("/user_signin", handlerSignInUser)
	http.HandleFunc("/user_add_info", handlerAddUserInfo)
	http.HandleFunc("/user_get_info", handlerGetInfoUser)
	http.HandleFunc("/company_signup", handlerSignUpCompany)
	http.HandleFunc("/company_signin", handlerSignInCompany)
	http.HandleFunc("/company_get_info", handlerGetInfoCompany)
	http.HandleFunc("/company_add_event", handlerAddEvent)
	http.HandleFunc("/get_all_events", handleGetAllEvents)
	http.HandleFunc("/get_company_events", handleGetCompanyEvents)
	http.HandleFunc("/user_join_event", handleUserJoinEvent)
	http.HandleFunc("/get_user_events", handleGetUserEvents)
	http.ListenAndServe(":8080", nil)
}