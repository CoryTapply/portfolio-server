package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
)

var localFullResVideoUploadPath = "./resources/fullResolution"
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

func saveFile(fileName string, originalFileName string, file io.Reader, videoId string, start string, end string) {
	newFile := writeFile(fileName, file)
	defer os.Remove(newFile.Name())

	var waitGroup sync.WaitGroup
	waitGroup.Add(3)

	go saveTrimmedVideo(newFile.Name(), start, end, localFullResVideoUploadPath+"/"+originalFileName, &waitGroup)
	go saveCompressVideo(newFile.Name(), start, end, localVideoUploadPath+"/"+videoId+".mp4", &waitGroup)
	go grabThumbnail(newFile.Name(), start, localThumbnailUploadPath+"/"+videoId+".jpg", &waitGroup)

	waitGroup.Wait()
	fmt.Println("All routines finished")
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

/**
 * Generates a Full resolution Video MP4 file trimmed to the desired length
 */
func saveTrimmedVideo(filePath string, start string, end string, outputFilePath string, waitGroup *sync.WaitGroup) {
	regex := regexp.MustCompile(`\\`)
	formattedFilePath := regex.ReplaceAllString(filePath, "/")

	optionsString := []string{
		"-ss", start,
		"-i", formattedFilePath,
		"-to", end,
		"-c", "copy",
		"-map", "0",
		outputFilePath, "-y",
	}
	runVideoCommand(optionsString)
	waitGroup.Done()
}

/**
 * Generates a Compressed Video MP4 file
 */
func saveCompressVideo(filePath string, start string, end string, outputFilePath string, waitGroup *sync.WaitGroup) {
	regex := regexp.MustCompile(`\\`)
	formattedFilePath := regex.ReplaceAllString(filePath, "/")

	// ffmpeg -i resources/uploaded/upload-814755955.mp4 -vf scale=-1:720 -c:v libx264 -crf 23 -preset medium -c:a copy -r 30 resources/uploaded/upload-814755955result.mp4

	optionsString := []string{
		"-ss", start,
		"-i", formattedFilePath,
		"-to", end,
		// "-preset", "veryslow",
		"-preset", "medium",
		"-vf", scale1080p,
		"-c:v", "libx264",
		"-crf", qualityHigh,
		"-filter_complex", "[0:a:1]volume=1.0[l];[0:a:0][l]amerge=inputs=2[a]",
		"-map", "0:v:0",
		"-map", "[a]",
		"-r", frameRate60,
		outputFilePath, "-y",
	}
	runVideoCommand(optionsString)
	waitGroup.Done()
}

/**
 * Generates a Thumbnail JPG file from the first frame of the video
 */
func grabThumbnail(filePath string, start string, outputFilePath string, waitGroup *sync.WaitGroup) {
	regex := regexp.MustCompile(`\\`)
	formattedFilePath := regex.ReplaceAllString(filePath, "/")
	optionsString := []string{"-ss", start, "-i", formattedFilePath, "-vframes", "1", "-vf", scale480p, "-q:v", qualityVeryHigh, outputFilePath, "-y"}
	runVideoCommand(optionsString)
	waitGroup.Done()
}

/**
 * Runs an FFMPEG command with the given options string array
 */
func runVideoCommand(optionsString []string) {
	cmd := exec.Command("ffmpeg", optionsString...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		fmt.Printf("cmd.Run() failed with %s\n", err)
		fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		return
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		fmt.Printf("Command Failed: %s\nCommand: %s\n", err, optionsString)
		fmt.Println("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
	} else {
		fmt.Println("----------------------------------------------")
		fmt.Printf("Successfully finished the \nCommand: %s\n", optionsString)
		fmt.Println("----------------------------------------------")
	}
}

/**
 * Method to get the video files in the given fileDirectory
 */
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
