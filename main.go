package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func heavy(ctx context.Context) {
	cpus := runtime.NumCPU()
	fmt.Println("CPUs: ", cpus)
	for i := 0; i < cpus; i++ {
		go func() {
			i := 0
			builder := strings.Builder{}
			for {
				select {
				case <-ctx.Done():
					fmt.Println("done")
					return
				default:
					builder.WriteString("a")
					hash := sha256.Sum256([]byte(builder.String()))
					i += int(hash[2])
				}
			}
		}()
	}
}

func main() {
	var ctx context.Context
	var cancel context.CancelFunc
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request to /")
		fmt.Fprintf(w, "Hello There, this is a new update")
	})

	r.HandleFunc("/envs", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request to /envs")
		log.Println("environment variables: ", os.Environ())
		//print all envs
		fmt.Fprintf(w, "environment variables: %v", os.Environ())
	})

	r.HandleFunc("/status/{status}", func(w http.ResponseWriter, r *http.Request) {
		status, err := strconv.Atoi(mux.Vars(r)["status"])
		if err != nil {
			log.Println("error converting status to int: ", err)
			fmt.Fprintf(w, "error converting status to int: %v", err)
		}
		log.Println("request to /status with status ", status)
		w.WriteHeader(status)
		fmt.Fprintf(w, "status %v", status)
	})

	r.HandleFunc("/testDb", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request to /testDb")
		//try to pint mysql, user is root and password is set in env as db_pass
		//the dns name of the db is an env var named db_host
		db, err := sql.Open("mysql", "root:"+os.Getenv("db_pass")+"@tcp("+os.Getenv("db_host")+":3306)/mysql")
		if err != nil {
			log.Println("error connecting to db: ", err)
			fmt.Fprintf(w, "error connecting to db: %v", err)
		}
		defer db.Close()

		if err := db.Ping(); err != nil {
			log.Println("error pinging db: ", err)
			fmt.Fprintf(w, "error pinging db: %v", err)
		}

		fmt.Fprintf(w, "connected to db")
	})

	r.HandleFunc("/selfDestruct", func(w http.ResponseWriter, r *http.Request) {
		//this will panic the app
		log.Println("request to /selfDestruct")
		fmt.Fprintf(w, "self destructing")
		panic("self destructing")
	})

	r.HandleFunc("/exit/{status}", func(w http.ResponseWriter, r *http.Request) {
		status := mux.Vars(r)["status"]
		log.Println("request to /exit with status ", status)
		fmt.Fprintf(w, "exiting with status %v", status)
		e, _ := strconv.Atoi(status)
		os.Exit(e)
	})
	r.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request to /hi")
		fmt.Fprintf(w, "hi, this is the request from an endpoint")
	})

	r.HandleFunc("/heavy", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request to /heavy")
		ctx, cancel = context.WithCancel(context.Background())
		go heavy(ctx)
		fmt.Fprintf(w, "heavy started")
	})
	r.HandleFunc("/heavy/stop", func(w http.ResponseWriter, r *http.Request) {
		log.Println("request to /heavy/stop")
		cancel()
		fmt.Fprintf(w, "stopped heavy")
	})
	r.HandleFunc("/sleep/{seconds}", func(w http.ResponseWriter, r *http.Request) {
		seconds := mux.Vars(r)["seconds"]
		log.Println("request to /sleep with seconds ", seconds)
		s, _ := strconv.Atoi(seconds)
		fmt.Fprintf(w, "sleeping for %v seconds\n", s)
		time.Sleep(time.Duration(s) * time.Second)
		fmt.Fprintf(w, "slept for %v seconds", s)
	})

	log.Println("starting on port 8080")
	log.Fatalln(http.ListenAndServe(":8080", r))
}
