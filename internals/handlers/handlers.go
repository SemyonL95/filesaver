package handlers

import (
	"crypto/rand"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
)

const FileChunk = 100000000 // 100MB

var fileNameRegex = regexp.MustCompile("^[a-z0-9_.@()-]+.txt$")

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

	err = a.FileStorage.Put(fileHeader.Filename, file)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		log.Printf("failed to save file: %v", err)

		return
	}

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

	file, err := a.FileStorage.Get(filename)
	if err, ok := err.(*os.PathError); ok {
		log.Printf("file not exists: %v", err)
		http.Error(w, "file not exist", http.StatusNotFound)

		return
	}

	if err != nil {
		log.Printf("error during get from storage: %v", err)
		http.Error(w, "error during get from", http.StatusInternalServerError)

		return
	}

	processID, err := randomProcessID()
	if err != nil {
		log.Printf("error generating processId: %v", err)
		http.Error(w, "error generating processID", http.StatusInternalServerError)

		return
	}

	cacheName := fmt.Sprintf("file:%s", processID)
	for word := range file {
		exists, err := a.Cache.Exists(r.Context(), cacheName, word)

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

		w.Write([]byte(word))
	}

	a.Cache.Remove(r.Context(), cacheName)

	return
}

func getFileExt(file multipart.File) (string, error) {
	fileHeader := make([]byte, 512) // http://golang.org/pkg/net/http/#DetectContentType

	if _, err := file.Read(fileHeader); err != nil {
		return "", err
	}

	if _, err := file.Seek(0, 0); err != nil {
		return "", err
	}

	contentType := http.DetectContentType(fileHeader)

	return contentType, nil
}

func randomProcessID() (s string, err error) {
	b := make([]byte, 8)
	_, err = rand.Read(b)
	if err != nil {
		return
	}
	s = fmt.Sprintf("%x", b)

	return
}

func validateFileName(fileName string) bool {
	return fileNameRegex.MatchString(fileName)
}
