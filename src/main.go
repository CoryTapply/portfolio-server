package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var indexEntryPoint = "./dist/index.html"

func main() {

	port, exists := os.LookupEnv("PORT")

	if !exists {
		port = "9090"
	}

	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/uploadVideo", uploadHandler).Methods("POST", "OPTIONS")
	api.HandleFunc("/getVideos", getVideosHandler).Methods("GET", "OPTIONS")

	// Serve static assets directly.
	r.PathPrefix("/resources/").Handler(http.FileServer(http.Dir(".")))
	r.PathPrefix("/dist/").Handler(http.FileServer(http.Dir(".")))

	// Health and Readiness Checks
	r.PathPrefix("/readiness_check").HandlerFunc(healthCheckHandler(indexEntryPoint))
	r.PathPrefix("/liveness_check").HandlerFunc(healthCheckHandler(indexEntryPoint))

	// Catch-all: Serve our JavaScript application's entry-point (index.html).
	r.PathPrefix("/").HandlerFunc(indexHandler(indexEntryPoint))

	http.Handle("/", r)

	fmt.Println("Server listening! On Port " + port)
	panic(http.ListenAndServe(":"+port, nil))
}
