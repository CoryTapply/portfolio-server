package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var indexEntryPoint = "./dist/index.html"

var port, _ = os.LookupEnv("PORT")
var env, _ = os.LookupEnv("ENVIRONMENT")

func main() {

	r := mux.NewRouter()

	// srv := &http.Server{
	// 	Addr: "0.0.0.0:" + port,
	// 	// Good practice to set timeouts to avoid Slowloris attacks.
	// 	WriteTimeout: time.Second * 300,
	// 	ReadTimeout:  time.Second * 600,
	// 	IdleTimeout:  time.Second * 600,
	// 	Handler:      r, // Pass our instance of gorilla/mux in.
	// }

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

	defer func() {
		if error := recover(); error != nil {
			fmt.Println("This is the error: ", error)
		}
	}()

	fmt.Println("Server listening! On Port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
	panic(err)
	// panic(srv.ListenAndServe())
}
