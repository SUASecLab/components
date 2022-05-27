package main

import (
	"database/sql"

	"github.com/SUASecLab/workadventure_admin_extensions/extensions"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"log"
	"net/http"
	"strconv"
)

func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	nr := r.URL.Query().Get("nr")
	components := r.URL.Query().Get("components")

	// find out whether user is admin
	isAdmin, errorMsg := extensions.UserIsAdmin(adminExtensions, uuid)
	if !isAdmin || len(errorMsg) != 0 {
		fmt.Fprintf(w, "You are not an administrator")
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
		return
	}

	err = db.Ping()
	if err != nil {
		log.Println("Could not ping DB:", err)
		return
	}

	var count int
	query := `SELECT COUNT(*) FROM components WHERE nr = ?`

	err = db.QueryRow(query, nrVal).Scan(&count)
	if err != nil {
		fmt.Fprintf(w, "Can not prepare insertion/update of data")
		log.Println("Can not query number of rows", err)
		return
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
