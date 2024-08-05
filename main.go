package main

import (
	"log"
	"log/slog"
	"net/http"
)

func main() {
	s := &Server{NewKeyValueStore()}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", s.GetKeys)
  mux.HandleFunc("GET /{key}", s.GetValue)
	mux.HandleFunc("PUT /{key}", s.PutOrCreateValue)
  mux.HandleFunc("DELETE /{key}", s.DeleteKeyValue)
	srv := &http.Server{
		Addr:    ":4000",
		Handler: mux,
	}

	slog.Info("Starting server on port 4000")

	log.Fatal(srv.ListenAndServe())
}
