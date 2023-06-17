package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v53/github"
	"io/ioutil"
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

	var response github.CheckSuiteEvent
	json.Unmarshal(body, &response)

	var installationId = response.GetInstallation().GetID()
	var commitSHA = response.GetCheckSuite().GetHeadSHA()

	if response.GetAction() == "requested" && response.GetCheckSuite() != nil {
		fmt.Println("success")

		ctx := context.Background()

		itr, _ := ghinstallation.NewKeyFromFile(http.DefaultTransport, 348670, installationId, "key2.pem")
		client := github.NewClient(&http.Client{Transport: itr})

		options := github.CreateCheckRunOptions{
			Name:    "Custom Test Visualizer",
			HeadSHA: commitSHA,
		}

		checkrun, _, _ := client.Checks.CreateCheckRun(ctx, "facilioo", "facilioo", options)

		var checkrunId = checkrun.GetID()

	} else {
		fmt.Println("fail")
	}
}
