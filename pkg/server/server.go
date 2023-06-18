package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
)

const DIR = "files/"

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hell0"))
}

func AddFilesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Add Endpoint Hit")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	formdata := r.MultipartForm
	files := formdata.File["file"]

	for i := range files {
		file, err := files[i].Open()
		if err != nil {
			http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
			return
		}
		defer file.Close()
		_, err = os.ReadFile(DIR + files[i].Filename)
		if err != nil {
			fmt.Println("File not found in store. Adding new file.")
		} else {
			http.Error(w, "File with same name present in store "+files[i].Filename, http.StatusBadRequest)
			return
		}
		f, err := os.Create(DIR + files[i].Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Files uploaded successfully : ")
		fmt.Fprintf(w, files[i].Filename+"\n")
	}
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("All Files in directory\n")
	files, err := ioutil.ReadDir(DIR)
	if err != nil {
		http.Error(w, "The id query parameter is missing", http.StatusBadRequest)
		return
	}
	var fileNames = make([]string, 0)
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileNames)
}

// for _, v := range m.File {
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Update....")
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return
	}
	// comaparing files ---
	// fileData, err := os.Open(DIR + handler.Filename)
	// if err != nil {
	// 	fmt.Println("File not found in store. Adding new file.")
	// } else {
	// 	buffer := make([]byte, 10240)
	// 	for {
	// 		_, err := fileData.Read(buffer)
	// 		if err != nil {
	// 			if err != io.EOF {
	// 				fmt.Println(err)
	// 			}
	// 			break
	// 		}
	// 	}
	// 	payloadBuffer := make([]byte, 10240)
	// 	for {
	// 		_, err := file.Read(payloadBuffer)
	// 		if err != nil {
	// 			if err != io.EOF {
	// 				fmt.Println(err)
	// 			}
	// 			break
	// 		}
	// 	}
	// 	if bytes.Equal(payloadBuffer, buffer) {
	// 		w.WriteHeader(http.StatusOK)
	// 		json.NewEncoder(w).Encode("FIle already present with same data")
	// 		return
	// 	}
	// 	fmt.Println("File not found in store. Adding new file.", fileData)
	// }
	// fileData.Close()
	// ----
	f, err := os.Create(DIR + handler.Filename) //, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	defer file.Close()
	w.Write([]byte("File updated successfully"))
}
func RemoveHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.FormValue("fileName")
	err := os.Remove(DIR + fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("FIle removed successfully")

}
func WcHandler(w http.ResponseWriter, r *http.Request) {
	count, err := wordCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(count)
}
func wordCount() (int, error) {
	var err error
	totalWords := 0
	files, err := os.ReadDir("files")
	if err != nil {
		return 0, err
	}
	re := regexp.MustCompile("[a-zA-Z']+")
	for _, file := range files {
		bs, err := os.ReadFile(DIR + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		text := string(bs)
		matches := re.FindAllString(text, -1)
		totalWords += len(matches)
	}
	fmt.Printf("Total number of words in all files are : %d \n", totalWords)
	return totalWords, nil
}
func FwHandler(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.FormValue("limit"))
	order := r.FormValue("order")
	resMap, err := wordsFreq(limit, order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resMap)
}
func wordsFreq(limit int, flag string) (map[string]int, error) {
	var err error
	files, err := os.ReadDir("files")
	re := regexp.MustCompile("[a-zA-Z']+")
	if err != nil {
		return nil, err
	}
	wordsCount := make(map[string]int)
	for _, file := range files {
		bs, err := os.ReadFile(DIR + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		text := string(bs)
		matches := re.FindAllString(text, -1)
		for _, match := range matches {
			wordsCount[match]++
		}
	}
	keys := make([]string, 0, len(wordsCount))
	for key := range wordsCount {
		keys = append(keys, key)
	}
	resMap := make(map[string]int)
	var countOfWords = len(wordsCount)
	if limit > -1 && limit < len(wordsCount) {
		countOfWords = limit
	}
	countOfWords = countOfWords - 1
	if flag == "asc" {
		sort.Slice(keys, func(i, j int) bool {
			return wordsCount[keys[i]] < wordsCount[keys[j]]
		})
		for idx, key := range keys {
			resMap[key] = wordsCount[key]
			if idx == countOfWords {
				break
			}
		}
	} else if flag == "dsc" {
		sort.Slice(keys, func(i, j int) bool {
			return wordsCount[keys[i]] > wordsCount[keys[j]]
		})
		for idx, key := range keys {
			resMap[key] = wordsCount[key]
			if idx == countOfWords {
				break
			}
		}
	} else {
		return nil, fmt.Errorf("valid flag are [asc, desc]")
	}

	return resMap, nil
}
