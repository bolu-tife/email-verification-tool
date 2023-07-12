package main

import "net/http"

type EmailVerifier struct {
	Email    string `default:"name@example.org"`
	Domain   string `default:"name@example.org"`
	UserName string
}

type EmailStatus struct {
	Email    string
	Domain   string
	UserName string
	Valid    bool
}
type Config struct {
	Port        string `default:"3000"`
	SenderEmail string
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string `json:"error"`
}
