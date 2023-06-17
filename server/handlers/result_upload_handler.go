package handlers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"bug-blaster-360/database"
	"bug-blaster-360/services"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v53/github"
)

func HandleResultUpload(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	db, _ := database.Connect("postgres://admin:bugblaster360@localhost:6543/mydb?sslmode=disable")

	tableName := "check_runs"
	fmt.Println(r.FormValue("commit"))
	data, _ := db.GetDataByField(tableName, "commit_hash", r.FormValue("commit"))

	convertedData := services.ConvertToString(data)
	installationId := convertedData[0]["installation_id"]
	installationId2 := convertedData[0]["check_run_id"]

	defer db.Close()

	r.ParseMultipartForm(10 << 20)

	files := r.MultipartForm.File["files"]
	//for _, file := range files {
	file := files[0]
	src, fileerr := file.Open()
	fmt.Println("src--------------------------------------------------------------------------", src)

	if fileerr != nil {
		fmt.Fprintln(w, fileerr)
		return
	}
	defer src.Close()

	//src2, _ := src.Read(*os)

	content, err := ioutil.ReadAll(src)

	// Convert the byte slice to a string
	fileContent := string(content)

	// Print the file content
	fmt.Println(fileContent)

	filename := "output.txt"

	erreee := WriteTextToFile(filename, fileContent)
	if erreee != nil {
		fmt.Println(err)
		return
	}

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

	value, _ := strconv.ParseInt(installationId, 10, 64)
	value2, _ := strconv.ParseInt(installationId2, 10, 64)

	ctx := context.Background()
	itr, _ := ghinstallation.NewKeyFromFile(http.DefaultTransport, 348670, value, "key2.pem")
	client := github.NewClient(&http.Client{Transport: itr})
	checkrun, _, err := client.Checks.UpdateCheckRun(ctx, "facilioo", "facilioo", value2, options)

	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Print(*checkrun.ID)

	filecontent, _ := ioutil.ReadAll(src)
	fmt.Print(string(filecontent))
	//}
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
	// if we os.Open returns an error then handle it

	file, _ := os.Open("output.txt")
	fmt.Println("opened file")

	// defer the closing of our xmlFile so that we can parse it later on
	defer file.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(file)

	// we initialize our TestRun entity
	var out TestRun
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'results' which we defined above
	xml.Unmarshal(byteValue, &out)

	return out
}

func testRunObjectToMarkdownTable(testRun TestRun) string {
	// parse to markdown string
	var markdown = "| Test2 | Status | Duration | Computer |\n" +
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

func WriteTextToFile(filename string, text string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(text)
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
