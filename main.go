package main

import (
	"git.aubm.net/kendo5731/http_files_client/app"
	"git.aubm.net/kendo5731/http_files_client/loggers"
	"git.aubm.net/kendo5731/http_files_client/remote"
)

func main() {
	app.SetGlobals()
	defer loggers.InitLoggers().Close()

	filesList, err := remote.GetFilesList()
	if err != nil {
		loggers.Error.Fatalf("failed to fetch new available files, %v", err)
	}

	ch := make(chan remote.Status)
	for _, files := range chunkFiles(filesList) {
		for _, file := range files {
			go remote.DownloadAndDelete(file, ch)
		}
		for range files {
			status := <-ch
			if !status.DownloadSuccess {
				loggers.Error.Printf("download failed [%v], elapsed time [%v], %v", status.File, status.Elapsed, status.Error)
			} else {
				loggers.Info.Printf("download success [%v], elapsed time [%v]", status.File, status.Elapsed)
				if !status.RemoveSuccess {
					loggers.Error.Printf("remove failed [%v], %v", status.File, status.Error)
				}
			}
		}
	}
}

func chunkFiles(files []string) [][]string {
	nbFiles := len(files)
	chunksSize := 5
	start := 0
	end := chunksSize
	var chunks [][]string
	for {
		if start <= nbFiles {
			if end > nbFiles {
				end = nbFiles
			}
			newchunk := files[start:end]
			if chunks == nil {
				chunks = [][]string{newchunk}
			} else {
				chunks = append(chunks, newchunk)
			}
			start += chunksSize
			end += chunksSize
		} else {
			return chunks
		}
	}
}
