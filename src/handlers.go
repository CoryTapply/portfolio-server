package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/rs/xid"
)

type ResponseData struct {
	Message string          `json:"message"`
	Videos  []VideoResponse `json:"videos"`
}

type VideoResponse struct {
	VideoLocation     string `json:"videoLocation"`
	ThumbnailLocation string `json:"thumbnailLocation"`
	ID                string `json:"id"`
	Title             string `json:"title"`
	Tags              string `json:"tags"`
	Game              string `json:"game"`
	HasVoice          bool   `json:"hasVoice"`
	ViewCount         int    `json:"viewCount"`
	Duration          string `json:"duration"`
}

type Video struct {
	ID        string
	Title     string
	Tags      string
	Game      string
	HasVoice  bool
	ViewCount int
	Duration  string
}

type UploadVideo struct {
	ID        string
	Title     string
	Tags      string
	Game      string
	HasVoice  bool
	ViewCount int
	Duration  string
}

func uploadHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
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

	videoID := xid.New()

	startTime := request.FormValue("start")
	endTime := request.FormValue("end")
	startTimeParsed, _ := strconv.ParseFloat(startTime, 64)
	endTimeParsed, _ := strconv.ParseFloat(endTime, 64)

	duration := fmt.Sprintf("%.2f", endTimeParsed-startTimeParsed)
	voiceChatEnabled, _ := strconv.ParseBool(request.FormValue("enableVoiceChat"))

	videoName := request.FormValue("videoName")

	createVideo(
		videoID.String(),
		videoName,
		request.FormValue("videoTags"),
		request.FormValue("videoGame"),
		voiceChatEnabled,
		duration,
	)

	// Write the file to disk
	fileEnding := getFileExtension(contentType)
	fileName := "upload-*" + fileEnding
	fullVideoName := trimFileEnding(handler.Filename) + "_" + videoName + ".mp4"
	saveFile(fileName, fullVideoName, file, videoID.String(), startTime, endTime)

	data := ResponseData{Message: "Successfully Uploaded File"}
	writer.WriteHeader(http.StatusOK)
	encodingError := json.NewEncoder(writer).Encode(data)
	if encodingError != nil {
		fmt.Println(encodingError)
	}
}

func getVideosHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if request.Method == "OPTIONS" {
		writer.WriteHeader(http.StatusOK)
		return
	}

	// rootDirectory := "."
	fileDirectory := "/resources/uploaded/"
	fileThumbnailDirectory := "/resources/thumbnails/"
	// files := getVideos(rootDirectory + fileDirectory)

	videos := getVideosFromDB()

	data := ResponseData{Message: "Successfully Loaded Files", Videos: generateVideoResponseObject(videos, request.Host, fileDirectory, fileThumbnailDirectory)}
	responseJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseJSON)
}

/**
 * Handler method for generating the json response object of the videos that were loaded from disk
 */
func generateVideoResponseObject(videos []Video, host string, directoryPath string, thumbnailDirectory string) (response []VideoResponse) {
	// Set the security scheme for the video urls
	scheme := "https://"
	if env == "LOCAL" {
		scheme = "http://"
	}

	for _, video := range videos {
		response = append(response, VideoResponse{
			VideoLocation:     scheme + host + directoryPath + video.ID + ".mp4",
			ThumbnailLocation: scheme + host + thumbnailDirectory + video.ID + ".jpg",
			ID:                video.ID,
			Title:             video.Title,
			Tags:              video.Tags,
			Game:              video.Game,
			HasVoice:          video.HasVoice,
			ViewCount:         video.ViewCount,
			Duration:          video.Duration,
		})
	}

	return response
}

/**
 * Handler method for just serving the index.html file
 */
func indexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}

	return http.HandlerFunc(fn)
}

/**
 * Handler method for responding to the health and readiness checks with a 204 code
 */
func healthCheckHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}

	return http.HandlerFunc(fn)
}

// func getResourceHandler(writer http.ResponseWriter, request *http.Request) {
// 	fmt.Println("Sending File...")
// 	fmt.Println(request.URL.Path)
// 	fmt.Println(formatRequest(request))
// 	http.ServeFile(writer, request, "."+request.URL.Path)
// }
