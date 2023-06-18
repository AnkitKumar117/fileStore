package cli

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func AllList() {
	resp, err := http.Get("http://localhost:8081/list/")
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func RemoveFile(filename []string) {
	queryParam := url.Values{}
	if filepath.Ext(filename[0]) != ".txt" {
		fmt.Println("Only .txt files are allowed")
		os.Exit(1)
	}
	queryParam.Add("fileName", filename[0])
	resp, err := http.Get("http://localhost:8081/remove?" + queryParam.Encode())
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func UpdateFile(filenames []string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fileName := filenames[0]
	if filepath.Ext(filenames[0]) != ".txt" {
		fmt.Println("Only .txt files are allowed")
		os.Exit(1)
	}
	file, err := os.Open(fileName)
	if err != nil {
		checkErr(err)
	}
	fw, err := w.CreateFormFile("file", fileName)
	if err != nil {
		checkErr(err)
	}
	if _, err := io.Copy(fw, file); err != nil {
		checkErr(err)
	}
	w.Close()
	defer file.Close()
	resp, err := http.Post("http://localhost:8081/update/", w.FormDataContentType(), &buf)
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}

func AddFiles(fileNames []string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	if len(fileNames) == 0 {
		fmt.Println("Please add .txt files.")
		os.Exit(1)
	}
	for i := range fileNames {
		if filepath.Ext(fileNames[i]) != ".txt" {
			fmt.Println("Only .txt files are allowed")
			os.Exit(1)
		}
		fileName := fileNames[i]
		file, err := os.Open(fileName)
		if err != nil {
			checkErr(err)
		}
		fw, err := w.CreateFormFile("file", fileName)
		if err != nil {
			checkErr(err)
		}
		if _, err := io.Copy(fw, file); err != nil {
			checkErr(err)
		}
		defer file.Close()
	}
	w.Close()
	resp, err := http.Post("http://localhost:8081/addFile/", w.FormDataContentType(), &buf)
	if err != nil {
		checkErr(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		checkErr(err)
	}
	fmt.Println(string(body))
}
func WordCount() {
	resp, err := http.Get("http://localhost:8081/wordCount/")
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Total count of words ---", string(body))
}

func FreqWords(limit int, order string) {
	queryParam := url.Values{}
	queryParam.Add("order", "dsc")
	queryParam.Add("limit", strconv.Itoa(10))
	if len(os.Args) > 2 {
		queryParam.Set("limit", strconv.Itoa(limit))
	}
	if len(os.Args) > 3 {
		queryParam.Set("order", order)
	}
	resp, err := http.Get("http://localhost:8081/freqWords?" + queryParam.Encode())
	if err != nil {
		fmt.Println(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
}
