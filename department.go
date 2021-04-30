// package department

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"/Users/clairerhoda/go/"

// 	_ "github.com/lib/pq"
// )

// type Department struct {
//     DepartmentId int `json:"department_id"`
// 	Name string `json:"name"`
// }

// func GETHandler(w http.ResponseWriter, r *http.Request) {
// 	db := main.OpenConnection()

// 	rows, err := db.Query("SELECT * FROM departments")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var departments []Department

// 	for rows.Next() {
// 		var department Department
// 		rows.Scan(&department.DepartmentId, &department.Name)
// 		departments = append(departments, department)
// 	}

// 	departmentsBytes, _ := json.MarshalIndent(departments, "", "\t")

// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write(departmentsBytes)

// 	defer rows.Close()
// 	defer db.Close()
// }

// func POSTHandler(w http.ResponseWriter, r *http.Request) {
//     db := OpenConnection()

//     var data []Department
// 	err := json.NewDecoder(r.Body).Decode(&data)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	sqlStatement := `INSERT INTO Departments (name) VALUES ($1)`

//     for i := range data {
//         _, err = db.Exec(sqlStatement, data[i].Name)
//         if err != nil {
//             w.WriteHeader(http.StatusBadRequest)
//             panic(err)
//         }
//     }
    
// 	w.WriteHeader(http.StatusOK)
// 	defer db.Close()
// }
