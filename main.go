package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Header().Set("Content-Type", "application/json")
	// Set a new cache-control header. If an existing "Cache-Control" header exists
	// it will be overwritten.
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	w.Write([]byte(`{"name": "Galbeyte"}`))
}

func snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func getSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for writing a new snippet..."))
}

func postSnippet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Save a new snippet..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home) // Restrict to match only with {$}
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/write", getSnippet)
	mux.HandleFunc("POST /snippet/write", postSnippet)

	log.Print("starting server on :4000")

	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}
