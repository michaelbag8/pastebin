package main

import (
	"encoding/json"
	"net/http"
	"time"
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
	var req struct {
    Title     string     `json:"title"`
    Content   string     `json:"content"`
    Language  string     `json:"language"`
    ExpiresAt *time.Time `json:"expires_at"`
    IsPublic  bool       `json:"is_public"`
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

	
	shortID := generateID(6)
	privatekey := ""
	if !req.IsPublic{
		privatekey = generateID(16)
	}

	result := Paste{
		Title: req.Title,
		Content: req.Content,
		Language: req.Language,
		Views: 0,
		CreatedAt:time.Now(),
		ExpiresAt: req.ExpiresAt,
		IsPublic:req.IsPublic ,
		PrivateKey: privatekey,
		ShortID: shortID,

	}

	mu.Lock()
	pasteStore[shortID]= result
	mu.Unlock()

		
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)

}
