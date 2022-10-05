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
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	message       string
	connectionUri string
)

func init() {
	message = os.Getenv("message")
	connectionUri = os.Getenv("connectionUri")
}

// returns a connection to the ipaas database
func connectToDB() error {
	//get context
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	//try to connect
	clientOptions := options.Client().ApplyURI(connectionUri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return fmt.Errorf("error pingin the database: %v", err)
	}

	return nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "wow i got this message in my environment: "+message)
	})

	r.HandleFunc("/conn", func(w http.ResponseWriter, r *http.Request) {
		//connect to mongo
		if err := connectToDB(); err != nil {
			fmt.Fprint(w, err)
			return
		}

		fmt.Fprint(w, "connected to the database")
	})

	log.Println("starting on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", r))
}
