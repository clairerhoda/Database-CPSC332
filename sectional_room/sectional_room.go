package sectional_room

import (
	"database/sql"
	"fmt"
	"encoding/json"
	"log"
	"net/http"

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

type SectionalRoom struct {
    SectionalRoomId int `json:"sectional_room_id"`
	Name string `json:"name"`
	QuantityOfRooms int `json:"quantity_of_rooms"`
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM sectional_rooms")
	if err != nil {
		log.Fatal(err)
	}

	var sectional_rooms []SectionalRoom

	for rows.Next() {
		var sectional_room SectionalRoom
		rows.Scan(&sectional_room.SectionalRoomId, &sectional_room.Name, &sectional_room.QuantityOfRooms)
		sectional_rooms = append(sectional_rooms, sectional_room)
	}

	sectionalRoomsBytes, _ := json.MarshalIndent(sectional_rooms, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(sectionalRoomsBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

    var data []SectionalRoom
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO sectional_rooms (name, quantity_of_rooms) VALUES ($1, $2)`
    for i := range data {
        _, err = db.Exec(sqlStatement, data[i].Name, data[i].QuantityOfRooms)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            panic(err)
        }
    }
    
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func SectionalRoomExecute() {
	http.HandleFunc("/getSectionalRooms", GETHandler)
	http.HandleFunc("/insertSectionalRooms", POSTHandler)
}