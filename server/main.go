package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"bug-blaster-360/database"

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

func connectDatabase() {
	connectionString := "postgres://username:password@host:port/database?sslmode=disable"
	db, err := database.Connect(connectionString)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	id := 1
	tableName := "your_table"
	err = db.InsertID(id, tableName)
	if err != nil {
		fmt.Println("Failed to insert ID:", err)
		return
	}
}

func handleGitHubRepoPost(w http.ResponseWriter, r *http.Request) {

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
