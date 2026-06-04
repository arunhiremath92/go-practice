package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestLoginHandler_Success(t *testing.T) {
	// Prepare userStore with a known hash for "alice"
	hashed, _ := hashPassword("pass@123")
	userStore["alice"] = hashed

	body := LoginRequest{
		Username: "alice",
		Password: "pass@123",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()

	LoginHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", resp.StatusCode)
	}
}

func TestLoginHandler_InvalidBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader([]byte("not-json")))
	w := httptest.NewRecorder()
	LoginHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid body, got %d", w.Code)
	}
}

func TestLoginHandler_MissingFields(t *testing.T) {
	body := LoginRequest{}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	LoginHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing fields, got %d", w.Code)
	}
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	body := LoginRequest{Username: "notfound", Password: "pass"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/login", bytes.NewReader(jsonBody))
	w := httptest.NewRecorder()
	LoginHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for user not found, got %d", w.Code)
	}
}

func TestTokenVerifiyHandler_Success(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/protect/alice/mocked-token", nil)
	req = mux.SetURLVars(req, map[string]string{
		"username": "alice",
		"token":    "mocked-token",
	})
	w := httptest.NewRecorder()
	TokenVerifiyHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for valid token, got %d", w.Code)
	}
}

func TestTokenVerifiyHandler_MissingParams(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/protect//", nil)
	req = mux.SetURLVars(req, map[string]string{
		"username": "",
		"token":    "",
	})
	w := httptest.NewRecorder()
	TokenVerifiyHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for missing params, got %d", w.Code)
	}
}
