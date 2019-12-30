package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

func getResource(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Sending File...")
	fmt.Println(request.URL.Path)
	fmt.Println(formatRequest(request))
	http.ServeFile(writer, request, "."+request.URL.Path)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello from the api!asdasdad")
}

func main() {
	r := mux.NewRouter()

	r.PathPrefix("/resources/").Handler(http.FileServer(http.Dir(".")))
	http.HandleFunc("/asd", handler)

	http.Handle("/", r)

	fmt.Println("Server listening!")
	fmt.Println("Server listeninglkjhsdkljhalsdasdasdasasdasddas")
	panic(http.ListenAndServe(":9090", nil))
}
