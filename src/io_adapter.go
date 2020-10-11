package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
)

var localUploadPath = "./resources/uploaded"

const (
	// Resolution
	scale720p  = "scale=-1:720"
	scale1080p = "scale=-1:1080"
	scale1440p = "scale=-1:1440"

	// Quality
	qualityHigh = "23"
	qualityLow  = "28"

	// Frame Rate
	frameRate30 = "30"
	frameRate60 = "60"
)

func saveFile(fileName string, file io.Reader, start string, end string) {
	newFile := writeFile(fileName, file)
	compressVideo(newFile.Name(), start, end, localUploadPath+"/testOutput6.mp4")
}

func writeFile(fileName string, file io.Reader) *os.File {

	tempFile, err := ioutil.TempFile(localUploadPath, fileName)
	if err != nil {
		os.Mkdir(localUploadPath, 0777)
		tempFile, err = ioutil.TempFile(localUploadPath, fileName)
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
	optionsString := []string{"-ss", start, "-i", formattedFilePath, "-to", end, "-vf", scale1080p, "-c:v", "libx264", "-crf", qualityLow, "-preset", "medium", "-c:a", "copy", "-r", frameRate30, outputFilePath, "-y"}
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
			err := os.Remove(filePath)
			if err != nil {
				fmt.Printf("os.Remove() failed with %s\n", err)
			}
		}
	}()
}
