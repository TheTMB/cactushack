package main

import (
	"net/http"
	"fmt"
	"database/sql"
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
	Id                  int
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
		var about string
		if st, err := checkDbErr(rows.Scan(&ID, &login, &password, &firstName, &lastName, &country, &city, &phone,
			&site_domain_address, &description, &about)); st {
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