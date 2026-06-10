package main

import (
	"log"
	"net/http"
	"sync"
	"fmt"
)

var (
	pasteStore = map[string]Paste{}
	mu         sync.RWMutex
)


func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	fmt.Println("Server starting on port 8080...")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}


	mux.HandleFunc("POST   /pastes     ", pasteHandler)
	mux.HandleFunc("GET    /pastes     ", pasteHandler)
	mux.HandleFunc("GET    /pastes/{id}", pasteHandler)
	mux.HandleFunc("DELETE /pastes/{id}", pasteHandler)
