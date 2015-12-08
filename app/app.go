package app

import "os"

// DestDir where local copies of downloaded files will be stored
var DestDir string

// Addr of the remote server, i.e. : 192.168.1.1:9999
var Addr string

// Token that must be passed as a query parameter in every requests to the remote
var Token string

// SetGlobals should be invoked during application startup to init
// global program variables
func SetGlobals() {
	if len(os.Args) < 4 {
		panic("Not enough arguments, correct usage is go run main.go /destination/dir 0.0.0.0:8888 mySecretToken")
	}
	DestDir = os.Args[1]
	Addr = os.Args[2]
	Token = os.Args[3]
}
