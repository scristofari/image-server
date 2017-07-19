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

func TestUploadImage(t *testing.T) {
	cases := []struct {
		maxSize int
		code    int
	}{
		{1 << 20, http.StatusCreated},    // MB
		{1 << 10, http.StatusBadRequest}, // KB
	}

	for _, c := range cases {
		uploadMaxSize = c.maxSize

		body, contentType, err := loadFormFile("../../golang.png")
		if err != nil {
			t.Error(err.Error())
		}

		r, _ := http.NewRequest("POST", "http://localhost/images/golang.png", body)
		r.Header.Set("Content-Type", contentType)
		w := httptest.NewRecorder()

		uploadHandler(w, r)

		if w.Code != c.code {
			t.Errorf("Expected %d, get %d", c.code, w.Code)
			t.Error(w.Body)
		}
	}
}

func TestGetImage(t *testing.T) {
	cases := []struct {
		in   string
		code int
	}{
		{"http://localhost/images/golang.png?r=50x0", http.StatusOK},
		{"http://localhost/images/golang.png?r=23x0&t=23x34", http.StatusBadRequest},
		{"http://localhost/images/golang.png", http.StatusBadRequest},
		{"http://localhost/images/golang.png?r=efkgjergx2", http.StatusBadRequest},
		{"http://localhost/images/golang.png?r=2500x4000", http.StatusBadRequest},
	}

	for _, c := range cases {
		r, err := http.NewRequest("GET", c.in, nil)
		if err != nil {
			t.Errorf(err.Error())
		}
		w := httptest.NewRecorder()

		imageHandler(w, r)

		if w.Code != c.code {
			t.Errorf("Expected %d, get %d", c.code, w.Code)
			t.Error(w.Body)
		}
	}
}

func loadFormFile(path string) (*bytes.Buffer, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, "", err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, "", err
	}

	body := new(bytes.Buffer)
	bw := multipart.NewWriter(body)

	fw, err := bw.CreateFormFile("image", fi.Name())
	if err != nil {
		return nil, "", err
	}
	fw.Write(fileContents)

	err = bw.Close()
	if err != nil {
		return nil, "", err
	}

	return body, bw.FormDataContentType(), nil
}
