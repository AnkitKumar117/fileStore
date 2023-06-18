package main

import (
	"net/http"
	"regexp"

	server "github.com/AnkitKumar117/fileStore/pkg/server"
)

var validPath = regexp.MustCompile("^/(addFile|list|remove|update|wordCount|freqWords)/")

func makeHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r)
	}
}

func main() {
	http.HandleFunc("/hello", server.HelloHandler)
	http.HandleFunc("/addFile/", makeHandler(server.AddFilesHandler))
	http.HandleFunc("/list/", makeHandler(server.ListHandler))
	http.HandleFunc("/update/", makeHandler(server.UpdateHandler))
	http.HandleFunc("/remove/", makeHandler(server.RemoveHandler))
	http.HandleFunc("/wordCount/", makeHandler(server.WcHandler))
	http.HandleFunc("/freqWords/", makeHandler(server.FwHandler))
	http.ListenAndServe(":8081", nil)
}
