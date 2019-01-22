package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"log"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// It's important that this is before your catch-all route ("/")
	api := r.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/users", getUsersHandler).Methods("GET")

	/*
		FileServer() is told the root of static files is "static".

		We want the URL to start with "/dist/".
		So if someone requests "/dist/example.txt",
		we want the server to send the file "/static/example.txt".

		In order to achieve this, we have to strip "/dist" from the URL,
		and the remaining will be the relative path compared to the root folder
		"static" which if we join gives:

		./static/example.txt

	*/
	r.PathPrefix("/dist/").Handler(http.StripPrefix("/dist/", http.FileServer(http.Dir("./static"))))

	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	// Serve static files for React App
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./build/static/"))))
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

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"id": "12345",
		"ts": time.Now().Format(time.RFC3339),
	}

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	w.Write(b)
}
