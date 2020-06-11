package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

var sampledata = [...]string{
	"Lorem",
	"ipsum",
	"dolor",
	"sit",
	"amet",
	"consectetur",
	"adipiscing",
	"elit.",
}

func Test_Upload(t *testing.T) {
	generateDummyFile("test.txt")
	file, err := os.Open("./test.txt")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	io.Copy(part, file)
	writer.Close()

	r, err := http.NewRequest("PUT", "http://localhost:8080/upload", body)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	response, err := client.Do(r)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		t.Log("Status code should to be 200")
		t.Fail()
	}
	os.Remove("./test.txt")
}

func Test_GetFile(t *testing.T) {
	generateDummyFile("../storage/test.txt")
	resp, err := http.Get("http://localhost:8080/files?filename=test.txt")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	words := []string{}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if len(words) != len(sampledata) {
		t.Log("sampledata size not equal to response")
		t.Fail()
	}
	os.Remove("../storage/test.txt")
}

func generateDummyFile(filename string) {

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	datawriter := bufio.NewWriter(file)

	counter := 0
	for i := 0; i < 100000; i++ {
		if counter == 8 {
			counter = 0
		}
		_, _ = datawriter.WriteString(sampledata[counter] + "\n")
		counter++
	}

	datawriter.Flush()
	file.Close()
}
