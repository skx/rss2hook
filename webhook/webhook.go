// webhook.go is a simple example program that will listen upon
// localhost:8080 and dump the contents of any HTTP POST received
// to the console.
//

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// HandleHook is called on any access to the server-root.
//
// If a POST request is received dump it to the console.  Regardless
// of the requested method we then send an "OK" response to the caller.
func HandleHook(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		content, _ := ioutil.ReadAll(r.Body)
		fmt.Printf("%s\n", content)
	}
	// Always return a response to the caller.
	io.WriteString(w, "OK\n")
}

func main() {

	// Bind our handler
	http.HandleFunc("/", HandleHook)

	// Launch the server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
