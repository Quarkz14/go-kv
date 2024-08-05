package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

const jsonContentType = "application/json"

type InMemoryStore interface {
	GetAllKeys() []string
	GetValue(key string) (any, error)
	PutOrCreateValue(key string, value any)
	DeleteValue(key string)
}

type AllKeys struct {
	Keys []string
}

type Put struct {
	Value string `json:"value"`
}
type Get struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type Server struct {
	store InMemoryStore
}

func (s *Server) GetKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	keysStored := s.store.GetAllKeys()
	if len(keysStored) == 0 {
		http.Error(w, "no keys found", http.StatusNotFound)
		return
	}
	j := &AllKeys{keysStored}
	err := json.NewEncoder(w).Encode(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) PutOrCreateValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)

	key := r.PathValue("key")
	var v Put

	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if v.Value == "" {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}

	s.store.PutOrCreateValue(key, v.Value)
}

func (s *Server) GetValue(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)

	key := r.PathValue("key")

	value, err := s.store.GetValue(key)
	if err != nil {
		if ok := errors.Is(err, KeyNotInMemory); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	kv := Get{
		Key:   key,
		Value: value,
	}

	err = json.NewEncoder(w).Encode(kv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteKeyValue(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")

	_, err := s.store.GetValue(key)
	if err != nil {
		if ok := errors.Is(err, KeyNotInMemory); ok {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	s.store.DeleteValue(key)
}
