package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"log"
	"net/http"
	"strconv"
)

func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query().Get("password")
	nr := r.URL.Query().Get("nr")
	components := r.URL.Query().Get("components")

	if password != editPassword {
		fmt.Fprintf(w, "Invalid password")
		return
	}

	nrVal, err := strconv.Atoi(nr)
	if err != nil {
		fmt.Fprintf(w, "Invalid workplace number")
		return
	}

	db, err := sql.Open("mysql", username+":"+password+"@("+hostname+":3306)/"+dbname+"?parseTime=true")
	defer db.Close()

	if err != nil {
		log.Println("Could not open DB connection:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("Could not ping DB:", err)
	}

	var count int
	query := `SELECT COUNT(*) FROM components WHERE nr = ?`

	err = db.QueryRow(query, nrVal).Scan(&count)
	if err != nil {
		fmt.Fprintf(w, "Can not prepare insertion/update of data")
		log.Println("Can not query number of rows", err)
	}

	if count == 0 {
		// Row does not exist -> can directly input
		query = `INSERT INTO components (nr, components) VALUES (?, ?)`
		_, err := db.Exec(query, nrVal, components)

		if err != nil {
			fmt.Fprintf(w, "Could not insert components")
			log.Println("Could not insert components", err)
		} else {
			fmt.Fprintf(w, "Inserted components")
		}
		return
	}
	// Rows exists -> update
	query = `UPDATE components SET components = ? WHERE nr = ?`
	_, err = db.Exec(query, components, nrVal)

	if err != nil {
		fmt.Fprintf(w, "Could not update components")
		log.Println("Could not update components:", err)
	} else {
		fmt.Fprintf(w, "Updated components")
	}
}
