package main

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"log"

	"./routes"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/ordishs/gocore"
)

var logger = gocore.Log("APIServer")

func main() {
	stats := gocore.Config().Stats()
	logger.Infof("STATS\n%s\nVERSION\n-------\n%s (%s)\n\n", stats, version, commit)

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan

		appCleanup()
		os.Exit(1)
	}()

	start()
}

func appCleanup() {
	logger.Infof("API Server shutting dowm...")
}

func start() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	api.HandleFunc("/books/{title}/page/{page:[0-9]+}", routes.GetBooks).Methods("GET")

	// Serve static assets directly.
	r.PathPrefix("/static/").Handler(http.FileServer(http.Dir("./build")))

	// Serve index page on all unhandled routes
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./build/index.html")
	})

	srv := &http.Server{
		Handler: handlers.LoggingHandler(os.Stdout, r),
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
