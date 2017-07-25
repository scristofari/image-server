package main

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/scristofari/image-server/part4/resizer"
)

var (
	server *httptest.Server
)

func init() {
	server = httptest.NewServer(Handlers())
}

func TestAccessToken(t *testing.T) {
	cases := []struct {
		user, password string
		code           int
	}{
		{"", "", http.StatusForbidden},
		{"app1", "", http.StatusForbidden},
		{"", "passApp1", http.StatusForbidden},
		{"foo", "passApp1", http.StatusForbidden},
		{"app1", "passApp1", http.StatusOK},
	}
	for _, c := range cases {
		r, err := http.NewRequest("GET", server.URL+"/access/token", nil)
		r.SetBasicAuth(c.user, c.password)
		if err != nil {
			t.Error(err)
		}
		res, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Error(err)
		}

		if res.StatusCode != c.code {
			t.Errorf("Expected %d, get %d", c.code, res.StatusCode)
			b, _ := ioutil.ReadAll(res.Body)
			t.Error(string(b))
		}
	}
}

func TestUploadImageAuth(t *testing.T) {
	cases := []struct {
		user, password string
		accessToken    bool
		code           int
	}{
		{"app1", "passApp1", true, http.StatusCreated},
		{"app1", "passApp1", false, http.StatusUnauthorized},
	}
	for _, c := range cases {
		r, err := http.NewRequest("GET", server.URL+"/access/token", nil)
		r.SetBasicAuth(c.user, c.password)
		if err != nil {
			t.Error(err)
		}
		res, err := http.DefaultClient.Do(r)
		if err != nil {
			t.Error(err)
		}
		defer res.Body.Close()
		urlLink, err := ioutil.ReadAll(res.Body)
		if err != nil {
			t.Error(err)
			return
		}

		body, contentType, err := loadFormFile("../../files/golang.png")
		if c.accessToken {
			r, err = http.NewRequest("POST", string(urlLink), body)
		} else {
			r, err = http.NewRequest("POST", server.URL+"/upload/token", body)
		}
		r.Header.Set("Content-Type", contentType)
		if err != nil {
			t.Error(err)
		}

		res, err = http.DefaultClient.Do(r)
		if err != nil {
			t.Error(err)
		}
		defer res.Body.Close()

		if res.StatusCode != c.code {
			t.Errorf("Expected %d, get %d", c.code, res.StatusCode)
			b, _ := ioutil.ReadAll(res.Body)
			t.Error(string(b))
		}
	}
}

func TestUploadImage(t *testing.T) {
	cases := []struct {
		maxSize int
		code    int
	}{
		{1 << 20, http.StatusCreated},    // MB
		{1 << 10, http.StatusBadRequest}, // KB
	}

	for _, c := range cases {
		resizer.UploadMaxSize = c.maxSize

		body, contentType, err := loadFormFile("../../files/golang.png")
		if err != nil {
			t.Error(err.Error())
		}

		r, _ := http.NewRequest("POST", "http://localhost/upload/csrf", body)
		r.Header.Set("Content-Type", contentType)
		w := httptest.NewRecorder()

		uploadHandleFunc(w, r)

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

		imageHandleFunc(w, r)

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
