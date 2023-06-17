package handlers

import (
	"bug-blaster-360/database"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v53/github"
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

	if response.GetAction() == "requested" && response.GetCheckSuite() != nil {
		var installationId = response.GetInstallation().GetID()
		var commitSHA = response.GetCheckSuite().GetHeadSHA()

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

		db, _ := database.Connect("postgres://admin:bugblaster360@localhost:6543/mydb?sslmode=disable")

		data := map[string]interface{}{
			"check_run_id":    checkrunId,
			"commit_hash":     commitSHA,
			"installation_id": installationId,
		}

		tableName := "check_runs"
		db.InsertData(tableName, data)

		db.Close()

		fmt.Println(checkrunId)

	} else {
		fmt.Println("fail")
	}
}
