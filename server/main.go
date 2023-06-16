package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v53/github"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/github/events", handleGitHubRepoPost)
	r.HandleFunc("/debug", handleGitHubAuthenticate)
	r.HandleFunc("/commits/{sha}/results", handleResultUpload).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", r))
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

func handleResultUpload(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sha, _ := vars["sha"]

	r.ParseMultipartForm(10 << 20)

	files := r.MultipartForm.File["results"]
	for _, file := range files {
		src, fileerr := file.Open()
		if fileerr != nil {
			fmt.Fprintln(w, fileerr)
			return
		}
		defer src.Close()
		filecontent, _ := ioutil.ReadAll(src)
		fmt.Print(string(filecontent))
	}

	fmt.Print(sha)
}

func handleGitHubAuthenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	itr, _ := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, 348670, "key.pem")
	client := github.NewClient(&http.Client{Transport: itr})

	app, _, err := client.Apps.Get(context.Background(), "")

	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Print(*app.Slug)
}
