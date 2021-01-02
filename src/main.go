package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var indexEntryPoint = "./dist/index.html"

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1/").Subrouter()
	api.HandleFunc("/uploadVideo", uploadHandler).Methods("POST")
	api.HandleFunc("/getVideos", getVideosHandler).Methods("GET")

	// Serve static assets directly.
	r.PathPrefix("/resources/").Handler(http.FileServer(http.Dir(".")))
	r.PathPrefix("/dist/").Handler(http.FileServer(http.Dir(".")))

	// Catch-all: Serve our JavaScript application's entry-point (index.html).
	r.PathPrefix("/").HandlerFunc(indexHandler(indexEntryPoint))

	http.Handle("/", r)

	fmt.Println("Server listening!")
	panic(http.ListenAndServe(":9090", nil))
}
