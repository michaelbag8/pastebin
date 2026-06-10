package main

import (
	"encoding/json"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func createPasteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req struct{
		Content string `json:"content"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err!=nil{
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == ""{
		http.Error(w, "Content is empty", http.StatusBadRequest)
		return
	}

	shortCode := generateID(6)
	if Paste.IsPublic == false{
		privateKey := generateID(16)
	}else{
		publickey := shortCode
	}

}
