package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	os.Exit(1)
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from the updated ipaas container :D v9.1 (sheesh version)")
	})

	log.Println("starting on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", r))
}
