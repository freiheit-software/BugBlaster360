package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"bug-blaster-360/database"
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
	fmt.Println("akk:", string(body))

	db, err := database.Connect("postgres://admin:bugblaster360@localhost:6543/mydb?sslmode=disable")

	data := map[string]interface{}{
		"check_run_id": id,
		"commit_hash":  string(body),
	}

	tableName := "check_runs"
	db.InsertData(tableName, data)

	defer db.Close()
	/* if err != nil {
		log.Fatal(err)
	} */

	fmt.Fprintln(w, "Received JSON data successfully!")
}
