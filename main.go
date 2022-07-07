package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hellooo, this server is gonna exit now")
		os.Exit(1)
	})

	log.Println("starting on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", r))
}
