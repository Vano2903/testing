package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from ipaas container :D")
	})

	log.Println("starting on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", r))
}
