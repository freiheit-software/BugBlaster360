package handlers

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"bug-blaster-360/models"
)

func Commit(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the path
	id := r.URL.Path[len("/commit/"):]
	fmt.Println("ID:", id)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Parse the JSON data
	var jsonData models.Commit
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return
	}

	// Access the ID and data fields from the parsed JSON
	fmt.Println("ID:", jsonData.ID)
	fmt.Println("Data:", jsonData.Data)

	// Handle the JSON data as required

	fmt.Fprintln(w, "Received JSON data successfully!")
}