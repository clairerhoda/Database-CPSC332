package department

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

type Department struct {
    DepartmentId int `json:"department_id"`
	Name string `json:"name"`
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	db := OpenConnection()

	rows, err := db.Query("SELECT * FROM departments")
	if err != nil {
		log.Fatal(err)
	}

	var departments []Department

	for rows.Next() {
		var department Department
		rows.Scan(&department.DepartmentId, &department.Name)
		departments = append(departments, department)
	}

	departmentsBytes, _ := json.MarshalIndent(departments, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(departmentsBytes)

	defer rows.Close()
	defer db.Close()
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
    db := OpenConnection()

    var data []Department
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sqlStatement := `INSERT INTO departments (name) VALUES ($1)`

    for i := range data {
        _, err = db.Exec(sqlStatement, data[i].Name)
        if err != nil {
            w.WriteHeader(http.StatusBadRequest)
            panic(err)
        }
    }
    
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}

func DepartmentExecute() {
	http.HandleFunc("/getDepartments", GETHandler)
	http.HandleFunc("/insertDepartments", POSTHandler)
}