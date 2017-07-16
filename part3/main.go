package main

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

const (
	mb = 1 << 20
)

var (
	outputDir = "tests/out"
)

func main() {
	http.Handle("/", handlers())

	log.Printf("Listening on port 8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/upload", uploadHandler).Methods("POST")
	r.HandleFunc("/images/{img}", imageHandler).Methods("GET")

	return r
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent from too large uploaded file
	r.Body = http.MaxBytesReader(w, r.Body, 5*mb)

	image, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get the image: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer image.Close()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f, err := os.OpenFile(outputDir+"/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, image)
	w.WriteHeader(http.StatusCreated)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	/** vars empty, in test case, we do not execute the router */
	hash := strings.Split(r.URL.Path, "/")
	filename := outputDir + "/" + hash[2]
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "failed to open file, "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q, err := getQuery(r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var i image.Image
	switch q.t {
	case "r":
		i = resize.Resize(q.s.width, q.s.height, img, resize.Lanczos2)
	case "t":
		i = resize.Thumbnail(q.s.width, q.s.height, img, resize.Lanczos2)
	}

	w.Header().Set("Cache-Control", "max-age=3600")
	png.Encode(w, i)
}

type query struct {
	t string
	s *size
}
type size struct {
	width  uint
	height uint
}

func getQuery(u *url.URL) (*query, error) {
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the query: %s", err)
	}

	if len(m) != 1 {
		return nil, fmt.Errorf("only one manipulation permitted")
	}

	if _, ok := m["r"]; ok {
		w, h, err := getSize(m["r"][0])
		return &query{
			t: "r",
			s: &size{
				width:  w,
				height: h,
			},
		}, err
	}

	if _, ok := m["t"]; ok {
		w, h, err := getSize(m["t"][0])
		return &query{
			t: "t",
			s: &size{
				width:  w,
				height: h,
			},
		}, err
	}

	return nil, fmt.Errorf("failed to get the manipulation %s", m)
}

func getSize(size string) (uint, uint, error) {
	hash := strings.Split(size, "x")
	if len(hash) != 2 {
		return 0, 0, fmt.Errorf("failed to get the proper size")
	}

	width, err := strconv.Atoi(hash[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get the proper width")
	}
	height, err := strconv.Atoi(hash[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get the proper height")
	}

	return uint(width), uint(height), nil
}
