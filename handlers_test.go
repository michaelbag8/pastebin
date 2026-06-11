package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func resetStore() {
	mu.Lock()
	defer mu.Unlock()

	pasteStore = make(map[string]Paste)
}

func TestCreatePasteHandler_PublicPaste(t *testing.T) {
	resetStore()

	body := `{
		"title":"test",
		"content":"hello world",
		"is_public":true
	}`

	req := httptest.NewRequest(http.MethodPost, "/pastes", strings.NewReader(body))
	rr := httptest.NewRecorder()

	createPasteHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, rr.Code)
	}

	var paste Paste
	if err := json.NewDecoder(rr.Body).Decode(&paste); err != nil {
		t.Fatal(err)
	}

	if paste.Content != "hello world" {
		t.Fatalf("expected content hello world got %s", paste.Content)
	}

	if paste.ShortID == "" {
		t.Fatal("expected short id to be generated")
	}
}

func TestCreatePasteHandler_PrivatePaste(t *testing.T) {
	resetStore()

	body := `{
		"title":"private",
		"content":"secret",
		"is_public":false
	}`

	req := httptest.NewRequest(http.MethodPost, "/pastes", strings.NewReader(body))
	rr := httptest.NewRecorder()

	createPasteHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected %d got %d", http.StatusCreated, rr.Code)
	}

	var paste Paste
	json.NewDecoder(rr.Body).Decode(&paste)

	if paste.PrivateKey == "" {
		t.Fatal("expected private key")
	}
}

func TestCreatePasteHandler_EmptyContent(t *testing.T) {
	resetStore()

	body := `{
		"title":"bad",
		"content":"",
		"is_public":true
	}`

	req := httptest.NewRequest(http.MethodPost, "/pastes", strings.NewReader(body))
	rr := httptest.NewRecorder()

	createPasteHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d", http.StatusBadRequest, rr.Code)
	}
}

func TestGetPasteHandler_NotFound(t *testing.T) {
	resetStore()

	req := httptest.NewRequest(http.MethodGet, "/pastes/abc", nil)
	req.SetPathValue("id", "abc")

	rr := httptest.NewRecorder()

	getPasteHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected %d got %d", http.StatusNotFound, rr.Code)
	}
}

func TestGetPasteHandler_PublicPaste(t *testing.T) {
	resetStore()

	pasteStore["abc"] = Paste{
		ShortID:  "abc",
		Content:  "hello",
		IsPublic: true,
	}

	req := httptest.NewRequest(http.MethodGet, "/pastes/abc", nil)
	req.SetPathValue("id", "abc")

	rr := httptest.NewRecorder()

	getPasteHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, rr.Code)
	}
}

func TestGetPasteHandler_PrivatePasteUnauthorized(t *testing.T) {
	resetStore()

	pasteStore["abc"] = Paste{
		ShortID:    "abc",
		Content:    "secret",
		IsPublic:   false,
		PrivateKey: "12345",
	}

	req := httptest.NewRequest(http.MethodGet, "/pastes/abc", nil)
	req.SetPathValue("id", "abc")

	rr := httptest.NewRecorder()

	getPasteHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestGetPasteHandler_PrivatePasteAuthorized(t *testing.T) {
	resetStore()

	pasteStore["abc"] = Paste{
		ShortID:    "abc",
		Content:    "secret",
		IsPublic:   false,
		PrivateKey: "12345",
	}

	req := httptest.NewRequest(http.MethodGet, "/pastes/abc", nil)
	req.SetPathValue("id", "abc")
	req.Header.Set("X-Private-Key", "12345")

	rr := httptest.NewRecorder()

	getPasteHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, rr.Code)
	}
}

func TestDeletePasteHandler_PublicPaste(t *testing.T) {
	resetStore()

	pasteStore["abc"] = Paste{
		ShortID:  "abc",
		IsPublic: true,
	}

	req := httptest.NewRequest(http.MethodDelete, "/pastes/abc", nil)
	req.SetPathValue("id", "abc")

	rr := httptest.NewRecorder()

	deletePasteHandler(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected %d got %d", http.StatusForbidden, rr.Code)
	}
}

func TestDeletePasteHandler_WrongKey(t *testing.T) {
	resetStore()

	pasteStore["abc"] = Paste{
		ShortID:    "abc",
		IsPublic:   false,
		PrivateKey: "secret",
	}

	req := httptest.NewRequest(http.MethodDelete, "/pastes/abc", nil)
	req.SetPathValue("id", "abc")
	req.Header.Set("X-Private-Key", "wrong")

	rr := httptest.NewRecorder()

	deletePasteHandler(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestDeletePasteHandler_CorrectKey(t *testing.T) {
	resetStore()

	pasteStore["abc"] = Paste{
		ShortID:    "abc",
		IsPublic:   false,
		PrivateKey: "secret",
	}

	req := httptest.NewRequest(http.MethodDelete, "/pastes/abc", nil)
	req.SetPathValue("id", "abc")
	req.Header.Set("X-Private-Key", "secret")

	rr := httptest.NewRecorder()

	deletePasteHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, rr.Code)
	}
}

func TestPublicPastes(t *testing.T) {
	resetStore()

	pasteStore["1"] = Paste{
		ShortID:  "1",
		IsPublic: true,
		Content:  "public",
	}

	pasteStore["2"] = Paste{
		ShortID:  "2",
		IsPublic: false,
		Content:  "private",
	}

	req := httptest.NewRequest(http.MethodGet, "/pastes", nil)
	rr := httptest.NewRecorder()

	publicPastes(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, rr.Code)
	}

	var pastes []Paste
	json.NewDecoder(rr.Body).Decode(&pastes)

	if len(pastes) != 1 {
		t.Fatalf("expected 1 public paste got %d", len(pastes))
	}
}

func TestPublicPastesEmpty(t *testing.T) {
	resetStore()

	req := httptest.NewRequest(http.MethodGet, "/pastes", nil)
	rr := httptest.NewRecorder()

	publicPastes(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected %d got %d", http.StatusOK, rr.Code)
	}

	var pastes []Paste
	json.NewDecoder(rr.Body).Decode(&pastes)

	if len(pastes) != 0 {
		t.Fatalf("expected empty slice got %d items", len(pastes))
	}
}
