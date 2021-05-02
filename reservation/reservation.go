package reservation

import (
	"database/sql"
	"fmt"
	"time"
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

type Reservation struct {
    ReservationId int `json:"reservation_id"`
	RoomId int `json:"room_id"`
	UserId string `json:"user_id"`
	PriorityLevel int `json:"priority_level"`
	AccountNumber int `json:"user_id"`
	StartTime string `json:"user_id"`
	EndTime string `json:"user_id"`
	Purpose string `json:"user_id"`
	NumberOfPeople int `json:"user_id"`
	CreatedAt string `json:"user_id"`
	CancelledPendingReopen string `json:"user_id"`
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM reservations")
	if err != nil {
		log.Fatal(err)
	}

	var reservations []Reservation

	for rows.Next() {
		var reservation Reservation
		rows.Scan(&reservation.ReservationId, &reservation.RoomId, &reservation.UserId, &reservation.PriorityLevel, &reservation.AccountNumber, &reservation.StartTime, &reservation.EndTime, &reservation.Purpose, &reservation.NumberOfPeople, &reservation.CreatedAt, &reservation.CancelledPendingReopen)
		reservations = append(reservations, reservation)
	}

	reservationsBytes, _ := json.MarshalIndent(reservations, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(reservationsBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	vars := mux.Vars(r)
    room_id := vars["room_id"]
    user_id := vars["user_id"]

	sqlStatement := `UPDATE rooms SET reservation_lock = NOT reservation_lock WHERE room_id = $1`
	_, err := db.Exec(sqlStatement, room_id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}


    var data []Reservation
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement = `INSERT INTO reservations (room_id, user_id, priority_level, account_number, start_time, end_time, purpose, number_of_people, created_at, cancelled_pending_reopen) VALUES ((SELECT room_id FROM rooms WHERE room_id='`+ room_id+ `'), (SELECT user_id FROM users WHERE user_id='`+ user_id+ `'), (SELECT access_level FROM users WHERE user_id='`+ user_id+ `'), $1, $2, $3, $4, $5, $6, $7)`
    for i := range data {
        _, err = db.Exec(sqlStatement, data[i].AccountNumber, time.Now(), time.Now(), data[i].Purpose, data[i].NumberOfPeople, time.Now(), time.Now())
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
    reservation_id := vars["reservation_id"]

	db := OpenConnection()

	sqlStatement := `DELETE FROM reservations WHERE reservation_id =  $1`
	_, err := db.Exec(sqlStatement, reservation_id)

	if err != nil {
		panic(err)
	  }
    
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func ReservationExecute(r *mux.Router) {
	subRouter := r.PathPrefix("/").Subrouter()
	http.HandleFunc("/getReservations", GETHandler)
	subRouter.HandleFunc("/insertReservations/{room_id}/{user_id}", POSTHandler).Methods("POST")
	subRouter.HandleFunc("/deleteReservation/{reservation_id}", DELETEHandler).Methods("DELETE")
}