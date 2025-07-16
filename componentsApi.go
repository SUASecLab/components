package main

import (
	"html"

	"github.com/SUASecLab/workadventure_admin_extensions/extensions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"fmt"
	"log"
	"net/http"
	"strconv"
)

func handleAPIRequest(w http.ResponseWriter, r *http.Request) {
	userToken := r.URL.Query().Get("token")
	nr := r.URL.Query().Get("nr")
	components := r.URL.Query().Get("components")

	// escape input
	userToken = html.EscapeString(userToken)
	nr = html.EscapeString(nr)
	components = html.EscapeString(components)

	// find out whether user is allowed to change the components
	decision, err := extensions.GetAuthDecision("http://" + sidecarUrl +
		"/auth?token=" + userToken + "&service=updateComponents")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorMsg := "Error while checking if user is allowed to update components"
		fmt.Fprintf(w, "%s", errorMsg)
		log.Println(errorMsg, err)
		return
	}

	if !decision.Allowed {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "You are not allowed to update the components")
		log.Println("Forbidden access attempt on request API")
		return
	}

	nrVal, err := strconv.Atoi(nr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Invalid workplace number")
		log.Println("Invalid workplace number provided")
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

	opts := options.Update().SetUpsert(true)
	filter := bson.D{{Key: "nr", Value: nrVal}}
	update := bson.D{{Key: "$set",
		Value: bson.D{{Key: "components", Value: components}}}}

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not update database:", err)
		fmt.Fprintln(w, "Could not update database")
		return
	}

	if result.UpsertedCount == 1 {
		fmt.Fprintln(w, "Inserted workplace description")
	} else if result.ModifiedCount == 1 {
		fmt.Fprintln(w, "Updated workplace description")
	} else {
		fmt.Fprintln(w, "Did not update workplace description. Are you sure the components changed?")
	}
}
