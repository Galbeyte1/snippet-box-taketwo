package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
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

func snippetWrite(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for writing a new snippet..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home) // Restrict to match only with {$}
	mux.HandleFunc("/snippet/view/{id}", snippetView)
	mux.HandleFunc("/snippet/write", snippetWrite)

	log.Print("starting server on :4000")

	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}
