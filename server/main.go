package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/github/events", handleGitHubRepoPost)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGitHubRepoPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Log the input data
	log.Printf("Received GitHub repo POST: %s", body)

	// Send a response
	fmt.Fprintln(w, "Received GitHub repo POST")
}
