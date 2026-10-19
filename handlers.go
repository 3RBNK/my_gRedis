package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type PutRequest struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

func (s *Server) healthHandler(w http.ResponseWriter) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "i'm alive <3",
	})
}

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	value, exp := s.kvs.Get(key)
	if !exp {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"key":   key,
		"value": value,
	})
}

func (s *Server) setHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	var req PutRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	if req.TTL < 0 {
		http.Error(w, "TTL must be positive", http.StatusBadRequest)
		return
	}

	ttl := time.Duration(req.TTL) * time.Second

	s.kvs.Put(key, req.Value, ttl)

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "saved",
	})
}

func (s *Server) delHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	s.kvs.Delete(key)

	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}
