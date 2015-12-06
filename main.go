package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

var destDir string
var hostPort string
var token string

func main() {
	initGlobals()
	defer initLogFile().Close()

	files, err := getFilesList()
	if err != nil {
		log.Fatal(err)
		return
	}

	ch := make(chan string)
	for _, file := range files {
		go downloadFile(file, ch)
	}
	for range files {
		fmt.Println(<-ch)
	}
}

func initGlobals() {
	if len(os.Args) < 4 {
		log.Fatal("Not enough arguments, correct usage is go run main.go /destination/dir 0.0.0.0:8888 mySecretToken")
	}
	destDir = os.Args[1]
	hostPort = os.Args[2]
	token = os.Args[3]
}

func initLogFile() *os.File {
	logResource, err := os.OpenFile("./errors.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logResource)
	return logResource
}

func getFilesList() ([]string, error) {
	resp, err := http.Get("http://" + hostPort + "/listFiles?token=" + token)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Request failed, server responded with a status code " + strconv.Itoa(resp.StatusCode))
	}

	respContent, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	var jsonResponse []string
	json.Unmarshal(respContent, &jsonResponse)

	return jsonResponse, nil
}

func downloadFile(filePathName string, ch chan<- string) {
	// create intermediate directories if needed
	dirname := filepath.Dir(filePathName)
	if dirname != "." {
		os.MkdirAll(destDir+"/"+dirname, os.ModePerm)
	}

	// create new resource
	out, err := os.Create(destDir + "/" + filePathName)
	defer out.Close()
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	// Perform a GET HTTP request to fetch the file
	requestQuery := "token=" + url.QueryEscape(token) + "&filename=" + url.QueryEscape(filePathName)
	requestURI := "http://" + hostPort + "/downloadFile?" + requestQuery
	resp, err := http.Get(requestURI)
	defer resp.Body.Close()

	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	if resp.StatusCode != 200 {
		ch <- fmt.Sprint("Get " + requestURI + ": " + resp.Status)
		return
	}

	// Copy the response body into the newly created file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	deleteRemoteFile(filePathName, ch)
}

func deleteRemoteFile(filePathName string, ch chan<- string) {
	// Build the DELETE HTTP request that delete the remote file from server
	requestQuery := "token=" + url.QueryEscape(token) + "&filename=" + url.QueryEscape(filePathName)
	requestURI := "http://" + hostPort + "/deleteFile?" + requestQuery
	req, err := http.NewRequest("DELETE", requestURI, nil)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	// Perform the request and handle response or error
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		ch <- fmt.Sprint("Delete " + requestURI + ": " + resp.Status)
		return
	}

	ch <- fmt.Sprintf("%v successfully downloaded from remote server.", filePathName)
}
