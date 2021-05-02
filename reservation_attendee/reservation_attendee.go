package reservation_attendee

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

type ReservationAttendee struct {
    ReservationId int `json:"reservation_id"`
	UserId int `json:"user_id"`
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	vars := mux.Vars(r)
    reservation_id := vars["reservation_id"]

	rows, err := db.Query(`SELECT reservations.reservation_id, users.user_id FROM reservations, users WHERE reservations.reservation_id = `+ reservation_id)
	if err != nil {
		log.Fatal(err)
	}

	var reservation_attendees []ReservationAttendee

	for rows.Next() {
		var reservation_attendee ReservationAttendee
		rows.Scan(&reservation_attendee.ReservationId, &reservation_attendee.UserId)
		reservation_attendees = append(reservation_attendees, reservation_attendee)
	}

	reservationAttendeesBytes, _ := json.MarshalIndent(reservation_attendees, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(reservationAttendeesBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	vars := mux.Vars(r)
    reservation_id := vars["reservation_id"]
    user_id := vars["user_id"]


    var data []ReservationAttendee
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO reservation_attendees (reservation_id, user_id) VALUES ($1, $2)`
	_, err = db.Exec(sqlStatement, reservation_id, user_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
    
    
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func ReservationAttendeeExecute(r *mux.Router) {
	subRouter := r.PathPrefix("/").Subrouter()
	subRouter.HandleFunc("/getReservationAttendees/{reservation_id}", GETHandler).Methods("GET")
	subRouter.HandleFunc("/insertReservationAttendees/{reservation_id}/{user_id}", GETHandler).Methods("POST")
}