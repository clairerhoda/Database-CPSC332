package room

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

type Room struct {
    RoomId int `json:"room _id"`
	BuildingId int `json:"building_id"`
	OfficeId int `json:"office_id"`
	SectionalRoomId int `json:"sectional_room_id"`
	Price float32 `json:"price"`
	RoomNumber int `json:"room_number"`
	Floor int `json:"floor"`
	Height int `json:"height"`
	Size int `json:"size"`
	Address string `json:"address"`
	Description string `json:"description"`
	GisCoordinate string `json:"gis_coordinate"`
	ReservationLock bool `json:"reservation_lock"`
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM rooms")
	if err != nil {
		log.Fatal(err)
	}

	var rooms []Room

	for rows.Next() {
		var room Room
		rows.Scan(&room.RoomId, &room.BuildingId, &room.OfficeId, &room.SectionalRoomId, &room.Price, &room.RoomNumber, &room.Floor, &room.Height, &room.Size, &room.Address, &room.Description, &room.GisCoordinate, &room.ReservationLock)
		rooms = append(rooms, room)
	}

	roomsBytes, _ := json.MarshalIndent(rooms, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(roomsBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    building_id := vars["building_id"]
    office_id := vars["office_id"]
    sectional_room_id := vars["sectional_room_id"]

	db := OpenConnection()

    var data []Room
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO rooms (building_id, office_id, sectional_room_id, price, room_number, floor, height, size, address, description, gis_coordinate, reservation_lock) VALUES ((SELECT building_id FROM buildings WHERE building_id='`+ building_id+ `'), (SELECT office_id FROM offices WHERE office_id='`+ office_id+ `'),(SELECT sectional_room_id FROM sectional_rooms WHERE sectional_room_id='`+ sectional_room_id + `'),$1, $2, $3, $4, $5, $6, $7, $8, $9)`
    for i := range data {
        _, err = db.Exec(sqlStatement, data[i].Price, data[i].RoomNumber, data[i].Floor, data[i].Height, data[i].Size, data[i].Address, data[i].Description, data[i].GisCoordinate, data[i].ReservationLock)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            panic(err)
        }
    }
    
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func RoomExecute(r *mux.Router) {
	subRouter := r.PathPrefix("/").Subrouter()
	http.HandleFunc("/getRooms", GETHandler)
	subRouter.HandleFunc("/insertRooms/{building_id}/{office_id}/{sectional_room_id}", POSTHandler).Methods("POST")
}