package handlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const FileChunk = 100000000 // 100MB

func (a *API) Upload(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s user-agent=%s", r.Method, r.Proto, r.URL.String(), r.UserAgent())
	if r.Method != http.MethodPut {
		http.Error(w, "Request method not allowed", http.StatusMethodNotAllowed)

		return
	}

	log.Printf("parsing multipart form")

	err := r.ParseMultipartForm(FileChunk)
	if err != nil {
		http.Error(w, "failed to parse file", http.StatusInternalServerError)
		log.Printf("failed to parse file: %v", err)

		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to read file form field", http.StatusInternalServerError)
		log.Printf("failed to read file form field: %v", err)

		return
	}
	defer file.Close()

	ext, err := getFileExt(file)
	if err != nil {
		http.Error(w, "cannot get file extension", http.StatusInternalServerError)
		log.Printf("cannot get file extension: %v", err)

		return
	}

	if ext != "text/plain; charset=utf-8" {
		errMsg := fmt.Sprintf("wrong file extension %s", ext)
		http.Error(w, errMsg, http.StatusInternalServerError)
		log.Print(errMsg)

		return
	}

	out, err := a.FileStorage.CreateEntry(fileHeader.Filename, file)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		log.Printf("failed to save file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "failed to write in file", http.StatusInternalServerError)
		log.Printf("failed to write in file: %v", err)

		return
	}

	w.Write([]byte("File saved"))

	return
}

func (a *API) GetFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s user-agent=%s", r.Method, r.Proto, r.URL.String(), r.UserAgent())

	if r.Method != http.MethodGet {
		http.Error(w, "Request method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Starts proccesing")

	filename := r.FormValue("filename")
	if valid := validateFileName(filename); !valid {
		log.Printf("filename query param required and have to be valid file name")
		http.Error(w, "filename query param required and have to be valid file name", http.StatusUnprocessableEntity)

		return
	}

	file, err := a.FileStorage.OpenEntry(filename)
	if err, ok := err.(*os.PathError); ok {
		log.Printf("file not exists: %v", err)
		http.Error(w, "file not exist", http.StatusNotFound)

		return
	}
	if err != nil {
		log.Printf("error during get from storage: %v", err)
		http.Error(w, "error during get from storage", http.StatusInternalServerError)

		return
	}
	defer file.Close()

	processID, err := randomProcessID()
	if err != nil {
		log.Printf("error generating processID: %v", err)
		http.Error(w, "error generating processID", http.StatusInternalServerError)

		return
	}

	cacheName := fmt.Sprintf("file:%s", processID)
	lineScanner := bufio.NewScanner(file)
	lineScanner.Split(bufio.ScanWords)
	for lineScanner.Scan() {
		word := lineScanner.Text()
		exists, err := a.Cache.Exists(r.Context(), cacheName, lineScanner.Text())

		if err != nil {
			log.Printf("redis error: %v", err)
			http.Error(w, "redis error", http.StatusInternalServerError)

			return
		}

		if exists {
			continue
		}

		_, err = a.Cache.Set(r.Context(), cacheName, word)
		if err != nil {
			log.Printf("redis error: %v", err)
			http.Error(w, "redis error", http.StatusInternalServerError)

			return
		}

		w.Write([]byte(fmt.Sprintf("%s\n", word)))
	}

	a.Cache.Remove(r.Context(), cacheName)

	if lineScanner.Err() != nil {
		log.Printf("scanner error: %v", err)
		http.Error(w, "scanner error", http.StatusInternalServerError)
	}

	return
}
