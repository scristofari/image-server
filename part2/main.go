package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/upload", uploadHandler).Methods("POST")
	r.HandleFunc("/images/{img}", imageHandler).Methods("GET")
	http.Handle("/", r)

	log.Printf("Listening on port 8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	defer file.Close()

	// Content-Length Mime-Version Content-type
	header.Header.Get("")

	f, err := os.OpenFile("./images/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)
	fmt.Println(file, header, err)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {

}
