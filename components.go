package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	r *mux.Router

	username string
	password string
	hostname string
	dbname   string

	adminExtensions string
)

func main() {
	log.SetFlags(0)

	username = os.Getenv("MYSQL_USER")
	password = os.Getenv("MYSQL_PASSWORD")
	hostname = os.Getenv("DB_HOSTNAME")
	dbname = os.Getenv("MYSQL_DATABASE")
	adminExtensions = os.Getenv("ADMIN_EXTENSIONS")

	r = mux.NewRouter()
	r.HandleFunc("/api/", handleAPIRequest)
	r.HandleFunc("/nr/{nr}", handleComponentsRequest)

	r.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		content, err := os.ReadFile("static/edit.html")

		if err != nil {
			fmt.Fprintf(w, "Can not serve file")
			return
		}
		fmt.Fprintf(w, string(content))
	})

	log.Println("Components is listening on port 8080")
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatalln("Components failed:", err)
	}
}
