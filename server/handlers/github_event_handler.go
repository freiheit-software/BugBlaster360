package handlers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

	log.Printf("Received GitHub repo POST: %s", body)

	fmt.Fprintln(w, "Received GitHub repo POST")

	// Parse Testfile
	var testRunObject = parseTestFile()
	var title = "Test Results"
	var summary = testRunObjectToSummary(testRunObject)
	var text = testRunObjectToMarkdownTable(testRunObject)
	var payload PatchCheckRunPayload
	var payloadOutput PatchCheckRunPayloadOutput
	payloadOutput.Title = title
	payloadOutput.Summary = summary
	payloadOutput.Text = text

	// Call Github Checks Api with formatted test results
	//decoder := json.NewDecoder(r.Body)
	//decoder.Decode(&payload)
	jsonPayload, _ := json.Marshal(payload)
	bytesPayload := bytes.NewReader(jsonPayload)
	// TODO: correct path
	req, err := http.NewRequest("PATCH", "http://???/repos/{owner}/{repo}/check-runs/{check_run_id}", bytesPayload)
	if err != nil {
		fmt.Println("Error is req: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// create a Client
	client := &http.Client{}
	// Do sends an HTTP request and
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in send req: ", err.Error())
		w.WriteHeader(400)
		//w.Write(err)
	}
	defer resp.Body.Close()
}

type PatchCheckRunPayload struct {
	Output PatchCheckRunPayloadOutput `json:"output"`
}

type PatchCheckRunPayloadOutput struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Text    string `json:"text"`
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
