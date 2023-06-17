package main

import (
	"fmt"
	"log"
	"net/http"
	//"os"

	"bug-blaster-360/dbmigrate"
	"bug-blaster-360/router"
)

func main() {
	router.SetupRouter()
	migrate()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func migrate() {
	fmt.Println("Logging a message to the console.")
	connectionString := "postgres://admin:bugblaster360@localhost:5432/mydb?sslmode=disable" //os.Getenv("DB_CONNECTION_STRING")
	fmt.Println("connectionString", connectionString)

	dbmigrate.ApplyMigrations()
}
