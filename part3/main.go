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
	uuid "github.com/satori/go.uuid"
)

const (
	mb            = 1 << 20
	presetMaxSize = 2000
)

var (
	outputDir     = "images"
	uploadMaxSize = 5 * mb
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
	// Prevent from too large uploaded file / PART 4
	r.Body = http.MaxBytesReader(w, r.Body, int64(uploadMaxSize))

	image, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get the image: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer image.Close()

	uuid := uuid.NewV4().String()
	f, err := os.OpenFile(outputDir+"/"+uuid+".png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, image)

	w.WriteHeader(http.StatusCreated) // Header status always before
	w.Write([]byte(fmt.Sprintf("%s://%s/images/%s", r.URL.Scheme, r.Host, uuid)))
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	/** vars from gorilla mux empty, in test case, we do not execute the router */
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
		i = resize.Resize(q.p.width, q.p.height, img, resize.Lanczos2)
	case "t":
		i = resize.Thumbnail(q.p.width, q.p.height, img, resize.Lanczos2)
	}

	w.Header().Set("Cache-Control", "max-age=3600")
	png.Encode(w, i)
}

type query struct {
	t string
	p *preset
}
type preset struct {
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

	for t, preset := range m {
		p, err := getPreset(preset[0])
		if err != nil {
			return nil, fmt.Errorf("failed to get the preset: %s", err)
		}

		return &query{
			t: t,
			p: p,
		}, nil
	}

	return nil, fmt.Errorf("failed to get the manipulation %s", m)
}

func getPreset(p string) (*preset, error) {
	hash := strings.Split(p, "x")
	if len(hash) != 2 {
		return nil, fmt.Errorf("failed to get the proper size")
	}

	width, err := strconv.Atoi(hash[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get the proper width")
	}
	height, err := strconv.Atoi(hash[1])
	if err != nil {
		return nil, fmt.Errorf("failed to get the proper height")
	}

	if height > presetMaxSize {
		return nil, fmt.Errorf("invalid preset, too high, got %d, max %d", height, presetMaxSize)
	}
	if width > presetMaxSize {
		return nil, fmt.Errorf("invalid preset, too high, got %d, max %d", width, presetMaxSize)
	}

	return &preset{
		width:  uint(width),
		height: uint(height),
	}, nil
}
