package main

import "net/http"

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Upload route"))
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8080", nil)
}
