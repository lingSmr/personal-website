package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	styleCss := NewServeFile("text/css", "style.css")

	http.Handle("/style.css", styleCss)
	http.HandleFunc("/video.mp4", VideoHandler)
	fs := http.FileServer(http.Dir("images"))
	http.Handle("/images/", http.StripPrefix("/images/", fs))

	// Обрабатываем корневой маршрут
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Запускаем сервер на порту 8080
	http.ListenAndServe(":8080", nil)
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
