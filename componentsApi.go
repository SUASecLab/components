package main

import (
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

	ctx, cancel, client, collection, success := connectToCollection(w)
	defer cancel()
	defer client.Disconnect(ctx)
	if !success {
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
	} else {
		if result.UpsertedCount == 1 {
			fmt.Fprintln(w, "Inserted workplace description")
		} else if result.ModifiedCount == 1 {
			fmt.Fprintln(w, "Updated workplace description")
		}
	}
}
