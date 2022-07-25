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

	adminExtensionsURL string
	externalToken      string
)

func main() {
	log.SetFlags(0)
	var exists bool

	username = os.Getenv("MYSQL_USER")
	password = os.Getenv("MYSQL_PASSWORD")
	hostname = os.Getenv("DB_HOSTNAME")
	dbname = os.Getenv("MYSQL_DATABASE")

	adminExtensionsURL, exists = os.LookupEnv("ADMIN_EXTENSIONS")
	if !exists {
		log.Fatalln("No admin extensions URL set")
	}

	externalToken, exists = os.LookupEnv("EXTERNAL_TOKEN")
	if !exists {
		log.Fatalln("No external token set")
	}

	r = mux.NewRouter()
	r.HandleFunc("/api/", handleAPIRequest)
	r.HandleFunc("/nr/{nr}", handleComponentsRequest)

	r.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		content, err := os.ReadFile("static/edit.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			msg := "Can not serve file"
			fmt.Fprintf(w, msg)
			log.Println(msg, err)
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
