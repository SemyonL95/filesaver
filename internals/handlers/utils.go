package handlers

import (
	"crypto/rand"
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
)

var fileNameRegex = regexp.MustCompile("^[a-z0-9_.@()-]+.txt$")

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
