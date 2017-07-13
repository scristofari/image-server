package main

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	server *httptest.Server
)

func init() {
	server = httptest.NewServer(handlers())
}

func TestUploadImage(t *testing.T) {
	file, err := os.Open("tests/golang.png")
	if err != nil {
		t.Error(err.Error())
	}
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		t.Error(err.Error())
	}
	fi, err := file.Stat()
	if err != nil {
		t.Error(err.Error())
	}

	body := new(bytes.Buffer)
	bw := multipart.NewWriter(body)

	fw, err := bw.CreateFormFile("image", fi.Name())
	if err != nil {
		t.Error("failed to add file")
	}
	fw.Write(fileContents)

	err = bw.Close()
	if err != nil {
		t.Error("failed to close the reader")
	}

	r, _ := http.NewRequest("POST", server.URL+"/upload", body)
	r.Header.Set("Content-Type", bw.FormDataContentType())

	w := httptest.NewRecorder()
	uploadHandler(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, get %d", http.StatusCreated, w.Code)
		t.Error(w.Body)
	}
}
