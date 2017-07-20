package main

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"

	"github.com/scristofari/image-server/part4/resizer"
)

var csrfHandler = csrf.Protect([]byte(securecookie.GenerateRandomKey(32)))

func main() {
	http.Handle("/", Handlers())

	log.Printf("Listening on port 8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Handlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/upload", authBasicHandleFunc(uploadHandleFunc)).Methods("POST")
	r.HandleFunc("/images/{img}", imageHandleFunc).Methods("GET")

	return r
}

func uploadHandleFunc(w http.ResponseWriter, r *http.Request) {
	// Prevent from too large uploaded file / PART 4
	r.Body = http.MaxBytesReader(w, r.Body, int64(resizer.UploadMaxSize))

	image, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get the image: %s", err.Error()), http.StatusBadRequest)
		return
	}
	defer image.Close()

	uuid, err := resizer.Uploadfile(image)

	w.WriteHeader(http.StatusCreated) // Header status always before
	w.Write([]byte(fmt.Sprintf("%s://%s/images/%s", r.URL.Scheme, r.Host, uuid)))
}

func imageHandleFunc(w http.ResponseWriter, r *http.Request) {
	/** vars from gorilla mux empty, in test case, we do not execute the router */
	hash := strings.Split(r.URL.Path, "/")

	q, err := resizer.GetQueryFromURL(r.URL)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get the query: %s", err.Error()), http.StatusBadRequest)
		return
	}

	i, err := resizer.Resize(hash[2], q)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to resize the image: %s", err.Error()), http.StatusBadRequest)
		return
	}

	w.Header().Set("Cache-Control", "max-age=3600")
	png.Encode(w, i)
}

func authBasicHandleFunc(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "failed to get auth basic credentials", http.StatusForbidden)
			return
		}

		err := resizer.CheckCredentials(user, pass)
		if err != nil {
			http.Error(w, "failed to sign in: "+err.Error(), http.StatusForbidden)
			return
		}

		f(w, r)
	}
}
