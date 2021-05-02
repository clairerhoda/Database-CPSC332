package user

import (
	"database/sql"
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "chese21"
	dbname   = "Rational_Room_Reservations"
)

func OpenConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

type User struct {
    UserId int `json:"user_id"`
	OfficeId int `json:"office_id"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Phone string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	AccessLevel int `json:"access_level"`
}


func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}

	var users []User

	for rows.Next() {
		var user User
		rows.Scan(&user.UserId, &user.OfficeId, &user.Email, &user.PasswordHash, &user.Phone, &user.FirstName, &user.LastName, &user.AccessLevel)
		users = append(users, user)
	}

	usersBytes, _ := json.MarshalIndent(users, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(usersBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    office_id := vars["office_id"]

	db := OpenConnection()

    var data []User
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO users (office_id, email, password_hash, phone, first_name, last_name, access_level) VALUES ((SELECT office_id FROM offices WHERE office_id='`+ office_id+ `'), $1, $2, $3, $4, $5, $6)`
    for i := range data {
        _, err = db.Exec(sqlStatement, data[i].Email, data[i].PasswordHash, data[i].Phone, data[i].FirstName, data[i].LastName, data[i].AccessLevel)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            panic(err)
        }
    }
    
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func DELETEHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    user_id := vars["user_id"]

	db := OpenConnection()

	sqlStatement := `DELETE FROM reservations WHERE user_id =  $1`
	_, err := db.Exec(sqlStatement, user_id)

	sqlStatement = `DELETE FROM users WHERE user_id =  $1`
	_, err = db.Exec(sqlStatement, user_id)

	if err != nil {
		panic(err)
	  }
    
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func UserExecute(r *mux.Router) {
	subRouter := r.PathPrefix("/").Subrouter()
	http.HandleFunc("/getUsers", GETHandler)
	subRouter.HandleFunc("/insertUsers/{office_id}", POSTHandler).Methods("POST")
	subRouter.HandleFunc("/deleteUser/{user_id}", DELETEHandler).Methods("DELETE")
}