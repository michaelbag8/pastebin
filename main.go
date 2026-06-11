package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	pasteStore = map[string]Paste{}
	mu         sync.RWMutex
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /pastes", createPasteHandler)
	mux.HandleFunc("GET /pastes/{id}", getPasteHandler)
	mux.HandleFunc("DELETE /pastes/{id}", deletePasteHandler)
	mux.HandleFunc("GET /pastes", publicPastes)

	fmt.Println("Server starting on port 8080...")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
