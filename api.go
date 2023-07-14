package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", makeHTTPHandleFunc(s.handleEmailVerificationCheck))

	log.Println("JSON API server running on port: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		log.Fatal("Error connecting to server")
	}
}

func (s *APIServer) handleEmailVerificationCheck(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	email := r.URL.Query().Get("email")

	emailVer, err := EmailVerificationProcess(email)

	emailStatus := emailVer.NewEmailStatus(err)

	return WriteJSON(w, http.StatusOK, emailStatus)
}
