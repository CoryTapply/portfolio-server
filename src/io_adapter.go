package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/rs/xid"
)

var localVideoUploadPath = "./resources/uploaded"
var localThumbnailUploadPath = "./resources/thumbnails"

const (
	// Resolution
	scale480p  = "scale=-1:480"
	scale720p  = "scale=-1:720"
	scale1080p = "scale=-1:1080"
	scale1440p = "scale=-1:1440"

	// Quality
	qualityBest     = "4"
	qualityVeryHigh = "12"
	qualityHigh     = "23"
	qualityLow      = "28"

	// Frame Rate
	frameRate30 = "30"
	frameRate60 = "60"
)

func saveFile(fileName string, file io.Reader, start string, end string) {
	newFile := writeFile(fileName, file)
	videoID := xid.New()
	// TODO: Save a trimmed copy at full resolution
	compressVideo(newFile.Name(), start, end, localVideoUploadPath+"/"+videoID.String()+".mp4")
	grabThumbnail(newFile.Name(), start, localThumbnailUploadPath+"/"+videoID.String()+".jpg")
}

func writeFile(fileName string, file io.Reader) *os.File {
	tempFile, err := ioutil.TempFile(localVideoUploadPath, fileName)
	if err != nil {
		os.Mkdir(localVideoUploadPath, 0777)
		tempFile, err = ioutil.TempFile(localVideoUploadPath, fileName)
		if err != nil {
			fmt.Println(err)
		}
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}

	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	return tempFile
}

func compressVideo(filePath string, start string, end string, outputFilePath string) {
	// "resources\uploaded\upload-814755955.mp4"

	regex := regexp.MustCompile(`\\`)
	formattedFilePath := regex.ReplaceAllString(filePath, "/")

	// ffmpeg -i resources/uploaded/upload-814755955.mp4 -vf scale=-1:720 -c:v libx264 -crf 23 -preset medium -c:a copy -r 30 resources/uploaded/upload-814755955result.mp4

	fmt.Println(formattedFilePath)
	optionsString := []string{"-ss", start, "-i", formattedFilePath, "-to", end, "-vf", scale1080p, "-c:v", "libx264", "-crf", qualityHigh, "-preset", "medium", "-c:a", "copy", "-map", "0", "-r", frameRate60, outputFilePath, "-y"}
	cmd := exec.Command("ffmpeg", optionsString...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}

	go func() {
		err := cmd.Wait()
		fmt.Println("YAY WE DID IT")
		if err != nil {
			fmt.Printf("cmd.Run() failed with %s\n", err)
		} else {
			fmt.Printf("No errors compressing video !! Wahoo\n")
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf("os.Remove() failed with %s\n", err)
			}
		}
	}()
}

func grabThumbnail(filePath string, start string, outputFilePath string) {
	regex := regexp.MustCompile(`\\`)
	formattedFilePath := regex.ReplaceAllString(filePath, "/")

	fmt.Println(formattedFilePath)
	optionsString := []string{"-ss", start, "-i", formattedFilePath, "-vframes", "1", "-vf", scale480p, "-q:v", qualityVeryHigh, outputFilePath, "-y"}
	cmd := exec.Command("ffmpeg", optionsString...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
	}

	go func() {
		err := cmd.Wait()
		fmt.Println("YAY WE DID IT")
		if err != nil {
			fmt.Printf("cmd.Run() failed with %s\n", err)
		} else {
			fmt.Printf("No errors with the thumbnail!! Wahoo\n")
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf("os.Remove() failed with %s\n", err)
			}
		}
	}()
}

func getVideos(fileDirectory string) (videos []os.FileInfo) {
	files, err := ioutil.ReadDir(fileDirectory)
	if err != nil {
		fmt.Printf("ioutil.ReadDir() failed with %s\n", err)
	}

	// Filter out any non-MP4 files
	for _, file := range files {
		if strings.Contains(file.Name(), ".mp4") {
			videos = append(videos, file)
		}
	}

	return videos
}
