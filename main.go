package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello motto"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	msg := fmt.Sprintf("Display a specific snippet with ID %d...", id)
	w.Write([]byte(msg))
}

func getSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for writing a new snippet..."))
}

func postSnippet(w http.ResponseWriter, r *http.Request) {
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
