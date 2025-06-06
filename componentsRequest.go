package main

import (
	"github.com/SUASecLab/workadventure_admin_extensions/extensions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func handleComponentsRequest(w http.ResponseWriter, r *http.Request) {
	// CORS and MIME type
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain")

	// Get workplace number
	variables := mux.Vars(r)
	nr := variables["nr"]

	// Get user ID
	userToken := r.URL.Query().Get("token")

	// Find out if user is allowed to receive workplace information
	decision, err := extensions.GetAuthDecision("http://" + sidecarUrl +
		"/auth?token=" + userToken + "&service=showComponents")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorMsg := "Could not receive authentication decision"
		log.Println(errorMsg, err)
		fmt.Fprintf(w, "%s", errorMsg)
		return
	}

	if !decision.Allowed {
		w.WriteHeader(http.StatusForbidden)
		errorMsg := "You are not allowed to access the components information"
		log.Println(errorMsg)
		fmt.Fprintf(w, "%s", errorMsg)
		return
	}

	// Get workplace description
	nrVal, err := strconv.Atoi(nr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid workplace number")
		log.Println("Could not convert workplace to number")
		return
	}

	ctx, cancel, client, collection, success := connectToCollection(w)
	defer cancel()
	defer client.Disconnect(ctx)
	if !success {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not connect to collection")
		return
	}

	var result bson.M

	err = collection.FindOne(ctx, bson.D{{Key: "nr", Value: nrVal}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintln(w, "No information stored")
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("Could not fetch document:", err)
			fmt.Fprintf(w, "An error occured while fetching the document")
		}
	} else {
		// We found the document
		fmt.Fprintf(w, "%v", result["components"])
	}
}
