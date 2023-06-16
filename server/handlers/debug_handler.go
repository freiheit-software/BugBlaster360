package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v53/github"
)

func HandleDebug(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := context.Background()

	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 348670, 38691856, "key2.pem")
	//itr, _ := ghinstallation.NewAppsTransportKeyFromFile(http.DefaultTransport, 348670, "key2.pem")
	client := github.NewClient(&http.Client{Transport: itr})

	/*
		options := github.CreateCheckRunOptions{
			Name:    "Custom Test Visualizer",
			HeadSHA: "bdfa19b89ffd80b29f07fb193d3bf60382159b1e",
		}
	*/

	var testRunObject = parseTestFile()
	var title = "Test Results"
	var summary = testRunObjectToSummary(testRunObject)
	var text = testRunObjectToMarkdownTable(testRunObject)

	status := "completed"
	conclusion := "success"

	options := github.UpdateCheckRunOptions{
		Name:       "Custom Test Visualizer",
		Status:     &status,
		Conclusion: &conclusion,
		Output: &github.CheckRunOutput{
			Title:   &title,
			Summary: &summary,
			Text:    &text,
		},
	}

	//checkrun, _, err := client.Checks.CreateCheckRun(ctx, "facilioo", "facilioo", options)
	checkrun, _, err := client.Checks.UpdateCheckRun(ctx, "facilioo", "facilioo", 14335840736, options)

	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Print(*checkrun.ID)

	/*
		app, _, err := client.Apps.Get(ctx, "")

		if err != nil {
			fmt.Print(err.Error())
		}

		fmt.Print(*app.Slug)
	*/
}

type TestRun struct {
	XMLName xml.Name `xml:"TestRun"`
	Results Results  `xml:"Results"`
}

type Results struct {
	XMLName         xml.Name         `xml:"Results"`
	UnitTestResults []UnitTestResult `xml:"UnitTestResult"`
}

type UnitTestResult struct {
	XMLName      xml.Name `xml:"UnitTestResult"`
	ExecutionId  string   `xml:"executionId,attr"`
	TestName     string   `xml:"testName,attr"`
	ComputerName string   `xml:"computerName,attr"`
	Duration     string   `xml:"duration,attr"`
	Status       string   `xml:"outcome,attr"`
}

func parseTestFile() TestRun {
	// Open our xmlFile
	xmlFile, err := os.Open("testfile.trx")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("opened file")

	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our TestRun entity
	var out TestRun
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'results' which we defined above
	xml.Unmarshal(byteValue, &out)

	return out
}

func testRunObjectToMarkdownTable(testRun TestRun) string {
	// parse to markdown string
	var markdown = "| Test | Status | Duration | Computer |\n" +
		"| ---- | ------ | -------- | -------- |\n"

	for i := 0; i < len(testRun.Results.UnitTestResults); i++ {
		var test = testRun.Results.UnitTestResults[i].TestName
		var status = testRun.Results.UnitTestResults[i].Status
		if status == "Passed" {
			status = "Passed :heavy_check_mark:"
		} else if status == "Failed" {
			status = "Failed :x:"
		}
		var duration = testRun.Results.UnitTestResults[i].Duration
		var computer = testRun.Results.UnitTestResults[i].ComputerName
		markdown = markdown + "| " + test + " | " + status + " | " + duration + " | " + computer + " |\n"
	}

	return markdown
}

func testRunObjectToSummary(testRun TestRun) string {
	var failingTests = 0
	for i := 0; i < len(testRun.Results.UnitTestResults); i++ {
		if testRun.Results.UnitTestResults[i].Status != "Passed" {
			failingTests++
		}
	}
	if failingTests == 0 {
		return "All tests passed! :rocket:"
	} else {
		return strconv.Itoa(failingTests) + " tests did not pass. :eyes:"
	}
}
