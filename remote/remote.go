package remote

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"git.aubm.net/kendo5731/http_files_client/app"
)

// GetFilesList returns a files of new available files
func GetFilesList() ([]string, error) {
	resp, err := http.Get(fmt.Sprintf("http://%v/listFiles?token=%v", app.Addr, app.Token))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with status %v", strconv.Itoa(resp.StatusCode))
	}

	respContent, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	var jsonResponse []string
	json.Unmarshal(respContent, &jsonResponse)

	return jsonResponse, nil
}

// DownloadAndDelete downloads a given file and remove it
// on the remote server. The function takes a channel in which it
// writes information about download.
func DownloadAndDelete(file string, ch chan<- Status) {
	var err error
	status := newStatus(file)
	startTime := time.Now()

	err = downloadFile(file)
	status.Elapsed = time.Now().Sub(startTime)
	if err == nil {
		err = deleteRemoteFile(file)
		if err != nil {
			status.RemoveSuccess = false
		}
	} else {
		status.DownloadSuccess = false
	}

	status.Error = err

	ch <- status
}

func downloadFile(file string) error {
	// create intermediate directories if needed
	dirname := filepath.Dir(file)
	if dirname != "." {
		os.MkdirAll(app.DestDir+"/"+dirname, os.ModePerm)
	}

	// create new resource
	out, err := os.Create(app.DestDir + "/" + file)
	defer out.Close()
	if err != nil {
		return err
	}

	// Perform a GET HTTP request to fetch the file
	requestURI := fmt.Sprintf("http://%v/downloadFile?token=%v&filename=%v", app.Addr, url.QueryEscape(app.Token), url.QueryEscape(file))
	resp, err := http.Get(requestURI)
	defer resp.Body.Close()

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("server responded with status %v\n", resp.Status)
	}

	// Copy the response body into the newly created file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func deleteRemoteFile(file string) error {
	// Build the DELETE HTTP request that delete the remote file from server
	requestURI := fmt.Sprintf("http://%v/deleteFile?token=%v&filename=%v", app.Addr, url.QueryEscape(app.Token), url.QueryEscape(file))
	req, err := http.NewRequest("DELETE", requestURI, nil)
	if err != nil {
		return err
	}

	// Perform the request and handle response or error
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("server responded with status %v", resp.Status)
	}

	return nil
}

// Status contains information about a file download
type Status struct {
	File            string
	DownloadSuccess bool
	RemoveSuccess   bool
	Error           error
	Elapsed         time.Duration
}

func newStatus(file string) Status {
	return Status{File: file, DownloadSuccess: true, RemoveSuccess: true, Error: nil}
}
