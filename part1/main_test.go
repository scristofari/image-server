package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpoad(t *testing.T) {
	r, _ := http.NewRequest("POST", "http://localhost/upload", nil)

	w := httptest.NewRecorder()
	uploadHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, get %d", http.StatusOK, w.Code)
		t.Error(w.Body)
	}
}
