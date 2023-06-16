package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v53/github"
)

func HandleDebug(w http.ResponseWriter, r *http.Request) {
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
