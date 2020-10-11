package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseData struct {
	Message string
}

func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Parse our multipart form, 10 << 32 specifies a maximum upload of 5 Gb
	request.ParseMultipartForm(10 << 32)
	file, handler, err := request.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	contentType := handler.Header.Get("Content-Type")

	// Write the file to disk
	fileEnding := getFileExtension(contentType)
	fileName := "upload-*" + fileEnding
	saveFile(fileName, file, request.FormValue("start"), request.FormValue("end"))

	data := ResponseData{Message: "Successfully Uploaded File"}
	writer.WriteHeader(http.StatusOK)
	encodingError := json.NewEncoder(writer).Encode(data)
	if encodingError != nil {
		fmt.Println(encodingError)
	}
}

// func getResourceHandler(writer http.ResponseWriter, request *http.Request) {
// 	fmt.Println("Sending File...")
// 	fmt.Println(request.URL.Path)
// 	fmt.Println(formatRequest(request))
// 	http.ServeFile(writer, request, "."+request.URL.Path)
// }
