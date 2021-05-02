package reserved_equipment

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

type ReservedEquipment struct {
    EquipmentId int `json:"equipment_id"`
	ReservationId int `json:"reservation_id"`
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()
	vars := mux.Vars(r)
    reservation_id := vars["reservation_id"]
	rows, err := db.Query(`SELECT equipment.equipment_id, reservations.reservation_id FROM equipment, reservations WHERE reservations.reservation_id= `+ reservation_id)
	if err != nil {
		log.Fatal(err)
	}

	var reserved_equipments []ReservedEquipment

	for rows.Next() {
		var reserved_equipment ReservedEquipment
		rows.Scan(&reserved_equipment.EquipmentId, &reserved_equipment.ReservationId)
		reserved_equipments = append(reserved_equipments, reserved_equipment)
	}

	reservedEquipmentBytes, _ := json.MarshalIndent(reserved_equipments, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(reservedEquipmentBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
		db := OpenConnection()

		vars := mux.Vars(r)
		equipment_id := vars["equipment_id"]
		reservation_id := vars["reservation_id"]
	
	
		var data []ReservedEquipment
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	
		sqlStatement := `INSERT INTO reserved_equipment (equipment_id, reservation_id) VALUES ($1, $2)`
		_, err = db.Exec(sqlStatement, equipment_id, reservation_id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}
		
		sqlStatement = `UPDATE equipment SET unavailable = NOT unavailable WHERE equipment_id = $1`
		_, err = db.Exec(sqlStatement, equipment_id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
		defer db.Close()
	}

func ReservedEquipmentExecute(r *mux.Router) {
	subRouter := r.PathPrefix("/").Subrouter()
	subRouter.HandleFunc("/getReservedEquipment/{reservation_id}", GETHandler).Methods("GET")
	subRouter.HandleFunc("/insertReservedEquipment/{equipment_id}/{reservation_id}", POSTHandler).Methods("POST")

}