package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello There, time to do some magic stuff, actually it's working :O")
	})

	
	r.HandleFunc("/envs", func(w http.ResponseWriter, r *http.Request) {
		//print all envs
		fmt.Fprintf(w, "environment variables: %v", os.Environ())
	})

	r.HandleFunc("/testDb", func(w http.ResponseWriter, r *http.Request) {
		//try to pint mysql, user is root and password is set in env as db_pass
		//the dns name of the db is an env var named db_host
		db, err := sql.Open("mysql", "root:"+os.Getenv("db_pass")+"@tcp("+os.Getenv("db_host")+":3306)/mysql")
		if err != nil {
			fmt.Fprintf(w, "error connecting to db: %v", err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			fmt.Fprintf(w, "error pinging db: %v", err)
		}

		fmt.Fprintf(w, "connected to db")
	})

	r.HandleFunc("/selfDestruct", func(w http.ResponseWriter, r *http.Request) {
		//this will panic the app
		fmt.Fprintf(w, "self destructing")
		panic("self destructing")
	})

	r.HandleFunc("/exit/{status}", func(w http.ResponseWriter, r *http.Request) {
		status := mux.Vars(r)["status"]
		fmt.Fprintf(w, "exiting with status %v", status)
		e, _ := strconv.Atoi(status)
		os.Exit(e)
	})
	r.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hi, this is the request from an endpoint")
	})
	log.Println("starting on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", r))
}
