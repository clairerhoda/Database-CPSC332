package building

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

type Building struct {
    BuildingId int `json:"building_id"`
	DepartmentId int `json:"department_id"`
	Name string `json:"name"`
	Address string `json:"address"`
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM buildings")
	if err != nil {
		log.Fatal(err)
	}

	var buildings []Building

	for rows.Next() {
		var building Building
		rows.Scan(&building.BuildingId, &building.DepartmentId, &building.Name, &building.Address)
		buildings = append(buildings, building)
	}

	buildingsBytes, _ := json.MarshalIndent(buildings, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(buildingsBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
    department_id := vars["department_id"]

	db := OpenConnection()

    var data []Building
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO buildings (department_id, name, address) VALUES ((SELECT department_id FROM departments WHERE department_id='`+ department_id+ `'), $1, $2)`
    for i := range data {
        _, err = db.Exec(sqlStatement, data[i].Name, data[i].Address)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            panic(err)
        }
    }
    
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func BuildingExecute(r *mux.Router) {
	subRouter := r.PathPrefix("/").Subrouter()
	http.HandleFunc("/getBuildings", GETHandler)
	subRouter.HandleFunc("/insertBuildings/{department_id}", POSTHandler).Methods("POST")
}