package equipment

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

type Equipment struct {
    EquipmentId int `json:"equipment_id"`
	RoomId int `json:"room_id"`
	Price float32 `json:"price"`
	Description string `json:"description"`
	Type string `json:"type"`
	Unavailable bool `json:"unavailable"`
	Fixed bool `json:"fixed"`
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM equipment")
	if err != nil {
		log.Fatal(err)
	}

	var equipments []Equipment

	for rows.Next() {
		var equipment Equipment
		rows.Scan(&equipment.EquipmentId, &equipment.RoomId, &equipment.Price, &equipment.Description, &equipment.Type, &equipment.Unavailable, &equipment.Fixed)
		equipments = append(equipments, equipment)
	}

	equipmentsBytes, _ := json.MarshalIndent(equipments, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(equipmentsBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    room_id := vars["room_id"]

	db := OpenConnection()

    var data []Equipment
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO equipment (room_id, price, description, type, unavailable, fixed) VALUES ((SELECT room_id FROM rooms WHERE room_id='`+ room_id+ `'), $1, $2, $3, $4, $5)`
    for i := range data {
        _, err = db.Exec(sqlStatement, data[i].Price, data[i].Description, data[i].Type, data[i].Unavailable, data[i].Fixed)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            panic(err)
        }
    }
    
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func EquipmentExecute(r *mux.Router) {
	subRouter := r.PathPrefix("/").Subrouter()
	http.HandleFunc("/getEquipment", GETHandler)
	subRouter.HandleFunc("/insertEquipment/{room_id}", POSTHandler).Methods("POST")
}