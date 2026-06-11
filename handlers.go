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
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is empty", http.StatusBadRequest)
		return
	}

	shortID := generateID(6)
	privatekey := ""
	if !req.IsPublic {
		privatekey = generateID(16)
	}

	result := Paste{
		Title:      req.Title,
		Content:    req.Content,
		Language:   req.Language,
		Views:      0,
		CreatedAt:  time.Now(),
		ExpiresAt:  req.ExpiresAt,
		IsPublic:   req.IsPublic,
		PrivateKey: privatekey,
		ShortID:    shortID,
	}

	mu.Lock()
	pasteStore[shortID] = result
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}

func getPasteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	mu.RLock()
	paste, ok := pasteStore[id]
	mu.RUnlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if paste.ExpiresAt != nil && time.Now().After(*paste.ExpiresAt) {
		http.Error(w, "paste has expired", http.StatusGone)
		return
	}
	if !paste.IsPublic {
		clientKey := r.Header.Get("X-Private-Key")
		if clientKey != paste.PrivateKey {
			http.Error(w, "unauthorised", http.StatusUnauthorized)
			return
		}
	}

	paste.Views++
	mu.Lock()
	pasteStore[id] = paste
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paste)
}

func deletePasteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	mu.RLock()
	paste, ok := pasteStore[id]
	mu.RUnlock()
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if paste.IsPublic {
		http.Error(w, "cannot delete public paste", http.StatusForbidden)
		return
	}

	clientKey := r.Header.Get("X-Private-Key")
	if clientKey != paste.PrivateKey {
		http.Error(w, "unauthorised", http.StatusUnauthorized)
		return
	}

	mu.Lock()
	delete(pasteStore, id)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "paste deleted"})
}

func publicPastes(w http.ResponseWriter, r *http.Request) {
	pastes := []Paste{}

	mu.RLock()
	for _, paste := range pasteStore {
		if !paste.IsPublic {
			continue
		}
		pastes = append(pastes, paste)
	}
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pastes)
}
