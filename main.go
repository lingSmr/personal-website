package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

const (
	port = ":8080"
)

func main() {
	const op = "main:main"
	log.Printf("%s:%s", op, "start initialization")

	styleCss := NewServeFile("text/css", "style.css")
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("images"))
	mux.Handle("/images/", http.StripPrefix("/images/", fs))
	mux.Handle("/style.css", styleCss)
	mux.HandleFunc("/video.mp4", VideoHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		const op = "handler: '/' "
		log.Printf("%s:%s", op, "someone visited main page")
		http.ServeFile(w, r, "index.html")
	})

	log.Printf("%s:%s", op, "Server started on port "+port)
	http.ListenAndServe(port, mux)
}

func VideoHandler(w http.ResponseWriter, r *http.Request) {
	vid, err := os.Open("images/video.mp4")
	if err != nil {
		http.Error(w, "could not open video", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, vid)
	if err != nil {
		http.Error(w, "could not send video", http.StatusInternalServerError)
		return
	}
}

type ServeFile struct {
	contentType string
	fileName    string
}

func NewServeFile(contentType, fileName string) *ServeFile {
	return &ServeFile{
		contentType: contentType,
		fileName:    fileName,
	}
}

func (sf *ServeFile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(sf.fileName)
	if err != nil {
		http.Error(w, "No such file", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Add("Content-Type", sf.contentType)
	io.Copy(w, file)
}
