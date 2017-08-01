package main

import (
	"log"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Upload route"))
}

func main() {
	r := http.NewServeMux()
	r.HandleFunc("/upload", uploadHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
