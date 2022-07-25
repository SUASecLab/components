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
	userToken := r.URL.Query().Get("token")
	nr := r.URL.Query().Get("nr")
	components := r.URL.Query().Get("components")

	// find out whether user is admin
	isAdmin, errorMsg := extensions.UserIsAdmin(adminExtensionsURL, userToken)
	if !isAdmin || len(errorMsg) != 0 {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "You are not an administrator")
		if len(errorMsg) > 0 {
			log.Println("Error while checking if user is admin:", errorMsg)
		}
		return
	}

	nrVal, err := strconv.Atoi(nr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid workplace number")
		return
	}

	db, err := sql.Open("mysql", username+":"+password+"@("+hostname+":3306)/"+dbname+"?parseTime=true")
	defer db.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Could not open DB connection"
		fmt.Fprintf(w, msg)
		log.Println(msg, err)
		return
	}

	err = db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Connection to database was lost"
		fmt.Fprintf(w, msg)
		log.Println(msg, err)
		return
	}

	var count int
	query := `SELECT COUNT(*) FROM components WHERE nr = ?`

	err = db.QueryRow(query, nrVal).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Can not prepare insertion/update of data"
		fmt.Fprintf(w, msg)
		log.Println(msg, err)
		return
	}

	if count == 0 {
		// Row does not exist -> can directly input
		query = `INSERT INTO components (nr, components) VALUES (?, ?)`
		_, err := db.Exec(query, nrVal, components)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			msg := "Could not insert components"
			fmt.Fprintf(w, msg)
			log.Println(msg, err)
		} else {
			fmt.Fprintf(w, "Inserted components")
		}
		return
	}
	// Rows exists -> update
	query = `UPDATE components SET components = ? WHERE nr = ?`
	_, err = db.Exec(query, components, nrVal)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Could not update components"
		fmt.Fprintf(w, msg)
		log.Println(msg, err)
	} else {
		fmt.Fprintf(w, "Updated components")
	}
}
