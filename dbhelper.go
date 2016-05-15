package main

import (
	"net/http"
	"fmt"
	"database/sql"
	"time"
	"strconv"
)

var db *sql.DB;

type user_st struct {
	ID               int
	Login            string
	Password         string
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

type company_st struct {
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

type company_event_st struct {
	ID          int
	Name        string
	Description string
	Start_event int64
	End_event   int64
	Token       int
}

type company_event_list_st map[string][]company_event_st

type company_user_event_st struct {
	ID       int
	user_id  int
	event_id int
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
		var firstName string
		var lastName string
		var country string
		var city string
		var university string
		var start_study int64
		var end_study int64
		var age int
		var work string
		var known_technology string
		var about string
		if st, err := checkDbErr(rows.Scan(&ID, &login, &password, &firstName, &lastName, &country, &city,
			&university, &start_study, &end_study, &age, &work, &known_technology, &about)); st {
			return nil, err
		}
		return user_st{ID, login, password, firstName, lastName, country, city, university, start_study, end_study, age, work, known_technology, about}, 0
	}
	return nil, 0
}

func GetUserByToken(token int) (interface{}, int) {
	rows, err := db.Query("SELECT * FROM Users WHERE id=?", token)
	if st, err := checkDbErr(err); st {
		return nil, err
	}
	if rows.Next() {
		var ID int
		var login string
		var password string
		var firstName string
		var lastName string
		var country string
		var city string
		var university string
		var start_study int64
		var end_study int64
		var age int
		var work string
		var known_technology string
		var about string
		if st, err := checkDbErr(rows.Scan(&ID, &login, &password, &firstName, &lastName, &country, &city,
			&university, &start_study, &end_study, &age, &work, &known_technology, &about)); st {
			return nil, err
		}
		return user_st{ID, login, password, firstName, lastName, country, city, university, start_study, end_study, age, work, known_technology, about}, 0
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

func UpdateUser(user user_st) (err int) {
	stmt, error := db.Prepare("update Users SET firstName=?, lastname=?, country=?, city=?, university=?, start_study=?, end_study=?," +
	"age=?, work=?, known_technology=?, about=? where id=?")
	if st, err := checkDbErr(error); st {
		return err;
	}

	if _, err := stmt.Exec(user.FirstName, user.LastName, user.Country, user.City, user.University, user.Start_study, user.End_study,
		user.Age, user.Work, user.Known_technology, user.About, user.ID); err != nil {
		if st, e := checkDbErr(err); st {
			return e
		}
	}
	return 0
}

func CreateCompany(company company_st) (err int) {
	stmt, error := db.Prepare("insert Companies SET login=?, password=?, firstName=?, lastName=?, country=?, city=?, " +
	"phone=?, site_domain_address=?, description=?")
	if st, err := checkDbErr(error); st {
		return err;
	}

	if _, err := stmt.Exec(company.Login, company.Password, company.FirstName, company.LastName, company.Country, company.City,
		company.Phone, company.Site_domain_address, company.Description); err != nil {
		if st, e := checkDbErr(err); st {
			return e
		}
	}
	return 0
}

func GetCompany(inputLogin string) (interface{}, int) {
	rows, err := db.Query("SELECT * FROM Companies WHERE login=?", inputLogin)
	if st, err := checkDbErr(err); st {
		return nil, err
	}
	if rows.Next() {
		var ID int
		var login string
		var password string
		var firstName string
		var lastName string
		var country string
		var city string
		var phone string
		var site_domain_address string
		var description string
		if st, err := checkDbErr(rows.Scan(&ID, &login, &password, &firstName, &lastName, &country, &city, &phone,
			&site_domain_address, &description)); st {
			return nil, err
		}
		return company_st{ID, login, password, firstName, lastName, country, city, phone, site_domain_address, description}, 0
	}
	return nil, 0
}

func GetCompanyByToken(token int) (interface{}, int) {
	rows, err := db.Query("SELECT * FROM Companies WHERE id=?", token)
	if st, err := checkDbErr(err); st {
		return nil, err
	}
	if rows.Next() {
		var ID int
		var login string
		var password string
		var firstName string
		var lastName string
		var country string
		var city string
		var phone string
		var site_domain_address string
		var description string
		if st, err := checkDbErr(rows.Scan(&ID, &login, &password, &firstName, &lastName, &country, &city, &phone,
			&site_domain_address, &description)); st {
			return nil, err
		}
		return company_st{ID, login, password, firstName, lastName, country, city, phone, site_domain_address, description}, 0
	}
	return nil, 0
}

func CreateEvent(event company_event_st) (err int) {
	stmt, error := db.Prepare("insert Events SET name=?, description=?, start_event=?, end_event=?, company_id=?");
	if st, err := checkDbErr(error); st {
		return err;
	}

	if _, err := stmt.Exec(event.Name, event.Description, event.Start_event, event.End_event, event.Token); err != nil {
		if st, e := checkDbErr(err); st {
			return e
		}
	}
	return 0
}

func GetEvents() (company_event_list_st, int) {
	rows, err := db.Query("SELECT * FROM Events")
	if st, err := checkDbErr(err); st {
		return nil, err
	}
	var results []company_event_st
	for rows.Next() {
		var ID int
		var name string
		var description string
		var start_event int64
		var end_event int64
		var company_id int
		if st, err := checkDbErr(rows.Scan(&ID, &name, &description, &start_event, &end_event, &company_id)); st {
			return nil, err
		}
		results = append(results, company_event_st{ID, name, description, start_event, end_event, company_id})
	}
	if len(results) > 0 {
		jsonFormatRes := company_event_list_st{
			"events":results,
		}
		return jsonFormatRes, 0
	}
	return nil, 0
}

func GetEventsByCompany(token int) (company_event_list_st, int) {
	rows, err := db.Query("SELECT * FROM Events where company_id=?", token)
	if st, err := checkDbErr(err); st {
		return nil, err
	}
	var results_old []company_event_st
	var results_new []company_event_st
	for rows.Next() {
		var ID int
		var name string
		var description string
		var start_event int64
		var end_event int64
		var company_id int
		if st, err := checkDbErr(rows.Scan(&ID, &name, &description, &start_event, &end_event, &company_id)); st {
			return nil, err
		}
		ev := company_event_st{ID, name, description, start_event, end_event, company_id}
		t := time.Now()
		format := t.Format("20060102150405")
		f, err := strconv.ParseInt(format, 10, 64)
		if err != nil {
			fmt.Printf("err convert to int %v", err)
		}
		if (f > end_event) {
			results_old = append(results_old, ev)
		} else {
			results_new = append(results_new, ev)
		}
	}
	if len(results_old) > 0 || len(results_new) > 0 {
		jsonFormatRes := company_event_list_st{
			"events_old":results_old,
			"events_new":results_new,
		}
		return jsonFormatRes, 0
	}
	return nil, 0
}

func CreateUserEvent(event company_user_event_st) (err int) {
	stmt, error := db.Prepare("insert UsersEvents SET event_id=?, user_id=?");
	if st, err := checkDbErr(error); st {
		return err;
	}

	if _, err := stmt.Exec(event.event_id, event.user_id); err != nil {
		if st, e := checkDbErr(err); st {
			return e
		}
	}
	return 0
}

func GetEventsByUser(token int) (company_event_list_st, int) {

	rows, err := db.Query("SELECT event_id FROM UsersEvents where user_id=?", token)
	if st, err := checkDbErr(err); st {
		return nil, err
	}
	var results_old []company_event_st
	var results_new []company_event_st
	for rows.Next() {
		var event_ID int
		if st, err := checkDbErr(rows.Scan(&event_ID)); st {
			return nil, err
		}

		rows, err := db.Query("SELECT * FROM Events where Id=?", event_ID)
		if st, err := checkDbErr(err); st {
			return nil, err
		}
		if rows.Next() {
			var ID int
			var name string
			var description string
			var start_event int64
			var end_event int64
			var company_id int
			if st, err := checkDbErr(rows.Scan(&ID, &name, &description, &start_event, &end_event, &company_id)); st {
				return nil, err
			}
			ev := company_event_st{ID, name, description, start_event, end_event, company_id}
			t := time.Now()
			format := t.Format("20060102150405")
			f, err := strconv.ParseInt(format, 10, 64)
			if err != nil {
				fmt.Printf("err convert to int %v", err)
			}
			if (f > end_event) {
				results_old = append(results_old, ev)
			} else {
				results_new = append(results_new, ev)
			}
		}
	}

	if len(results_old) > 0 || len(results_new) > 0 {
		jsonFormatRes := company_event_list_st{
			"events_old":results_old,
			"events_new":results_new,
		}
		return jsonFormatRes, 0
	}
	return nil, 0
}
