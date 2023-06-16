package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleGitHubRepoPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received GitHub repo POST: %s", body)

	fmt.Fprintln(w, "Received GitHub repo POST")
}
