package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func HandleResultUpload(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(10 << 20)

	files := r.MultipartForm.File["results"]
	for _, file := range files {
		src, fileerr := file.Open()
		if fileerr != nil {
			fmt.Fprintln(w, fileerr)
			return
		}
		defer src.Close()
		filecontent, _ := ioutil.ReadAll(src)
		fmt.Print(string(filecontent))
	}
}
