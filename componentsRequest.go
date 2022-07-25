package main

import (
	"database/sql"

	"github.com/SUASecLab/workadventure_admin_extensions/extensions"
	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func handleComponentsRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get workplace number
	variables := mux.Vars(r)
	nr := variables["nr"]

	// Get user ID
	userToken := r.URL.Query().Get("token")

	// Validate uuid
	exists, errorMsg := extensions.UserExists(adminExtensionsURL, userToken)
	if !exists {
		w.WriteHeader(http.StatusForbidden)
		log.Println(errorMsg)
		fmt.Fprintf(w, errorMsg)
		return
	}

	nrVal, err := strconv.Atoi(nr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
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

	var components string

	query := `SELECT components FROM components WHERE nr = ?`
	err = db.QueryRow(query, nrVal).Scan(&components)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := "Could not query database"
		fmt.Fprintf(w, msg)
		log.Println(msg, err)
		return
	}

	fmt.Fprintf(w, components)
}
