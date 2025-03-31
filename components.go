package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	r *mux.Router

	username string
	password string
	hostname string
	dbname   string

	sidecarUrl string
)

func connectToCollection(w http.ResponseWriter) (context.Context, context.CancelFunc, *mongo.Client, *mongo.Collection, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://"+username+":"+
		password+"@"+hostname+":27017/"+dbname))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not connect to database")
		return ctx, cancel, client, nil, false
	}

	return ctx, cancel, client, client.Database(dbname).Collection("components"), true
}

func main() {
	log.SetFlags(0)
	var exists bool

	username = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	hostname = os.Getenv("DB_HOST")
	dbname = os.Getenv("DB_NAME")

	sidecarUrl, exists = os.LookupEnv("SIDECAR_URL")
	if !exists {
		log.Fatalln("No sidecar URL set")
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
			fmt.Fprintf(w, "%s", msg)
			log.Println(msg, err)
			return
		}
		fmt.Fprintf(w, "%s", string(content))
	})

	log.Println("Components is listening on port 8080")
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatalln("Components failed:", err)
	}
}
