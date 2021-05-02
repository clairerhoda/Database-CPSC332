package user_proxy

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

type UserProxy struct {
    UserProxyId int `json:"user_proxy_id"`
	ParentUser int `json:"parent_user"`
	AccessLevel int `json:"access_level"`
	PasswordHash string `json:"password_hash"`
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM user_proxies")
	if err != nil {
		log.Fatal(err)
	}

	var user_proxies []UserProxy

	for rows.Next() {
		var user_proxy UserProxy
		rows.Scan(&user_proxy.UserProxyId, &user_proxy.ParentUser, &user_proxy.AccessLevel, &user_proxy.PasswordHash)
		user_proxies = append(user_proxies, user_proxy)
	}

	userProxiesBytes, _ := json.MarshalIndent(user_proxies, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(userProxiesBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    parent_user := vars["parent_user"]

	db := OpenConnection()

    var data []UserProxy
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO user_proxies (parent_user, access_level, password_hash) VALUES ((SELECT user_id FROM users WHERE user_id='`+ parent_user + `'), $1, $2)`
    for i := range data {
        _, err = db.Exec(sqlStatement, data[i].AccessLevel, data[i].PasswordHash)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            panic(err)
        }
    }
    
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func UserProxyExecute(r *mux.Router) {
	subRouter := r.PathPrefix("/").Subrouter()
	http.HandleFunc("/getUserProxies", GETHandler)
	subRouter.HandleFunc("/insertUserProxies/{parent_user}", POSTHandler).Methods("POST")
}