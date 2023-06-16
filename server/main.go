package main

import (
	"fmt"
	"log"
	"net/http"

	"bug-blaster-360/database"
	"bug-blaster-360/router"
)

func main() {
	router.SetupRouter()
	log.Fatal(http.ListenAndServe(":8080", nil))
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
