package router

import (
	"net/http"

	"bug-blaster-360/handlers"
)

func SetupRouter() {
	http.HandleFunc("/health", handlers.Health)
	http.HandleFunc("/commit/", handlers.Commit)
	http.HandleFunc("/github/events", handlers.HandleGitHubRepoPost)
}

