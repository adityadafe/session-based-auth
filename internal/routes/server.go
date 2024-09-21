package routes

import (
	"fmt"
	"net/http"
	"session-based-auth/internal/db"
)

type APIServer struct {
	listenAddr string
	store      db.Storage
}

func NewApiServer(listenAddr string, store db.Storage) *APIServer {
	return &APIServer{listenAddr: listenAddr, store: store}
}

func (s *APIServer) Run() {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /sign-up", makeHTTPHandlerFunc(s.signUpHandler))
	mux.HandleFunc("POST /sign-in", makeHTTPHandlerFunc(s.signInHandler))
	mux.HandleFunc("POST /sign-out", makeHTTPHandlerFunc(s.signOutHandler))
	mux.HandleFunc("GET /protected-route", makeHTTPHandlerFunc(s.signOutHandler))

	fmt.Println("Server is listening on http://localhost", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}
