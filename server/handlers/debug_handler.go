package handlers

import (
	"net/http"
)

func HandleDebug(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	/* ctx := context.Background()

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 348670, 38691856, "key2.pem")
	//itr, _ := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, 348670, "key2.pem")
	client := github.NewClient(&http.Client{Transport: itr}) */

	/*
		options := github.CreateCheckRunOptions{
			Name:    "Custom Test Visualizer",
			HeadSHA: "bdfa19b89ffd80b29f07fb193d3bf60382159b1e",
		}
	*/

	/*
		app, _, err := client.Apps.Get(ctx, "")

		if err != nil {
			fmt.Print(err.Error())
		}

		fmt.Print(*app.Slug)
	*/
}
