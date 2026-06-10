package main

import (
	"encoding/json"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err:= json.NewEncoder(w).Encode(map[string]string{
		"status":"ok",
	})
	if err!=nil{
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
