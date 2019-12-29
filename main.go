package main

import (
	"fmt"
	"net/http"
)

func getResource(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Sending File...")
	fmt.Println(request.URL.Path)
	http.ServeFile(writer, request, "."+request.URL.Path)
}

func main() {
	http.HandleFunc("/", getResource)
	panic(http.ListenAndServe("localhost:9090", nil))
}
