package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func handleComponentsRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	variables := mux.Vars(r)
	nr := variables["nr"]

	nrVal, err := strconv.Atoi(nr)
	if err != nil {
		fmt.Fprintf(w, "Invalid number")
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

	var components string

	query := `SELECT components FROM components WHERE nr = ?`
	err = db.QueryRow(query, nrVal).Scan(&components)

	if err != nil {
		log.Println("Could not query DB", err)
		return
	}

	fmt.Fprintf(w, components)
}
