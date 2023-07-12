package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func GetConfig() *Config {
	loadEnv()

	return &Config{
		Port:        os.Getenv("PORT"),
		SenderEmail: os.Getenv("SENDER_EMAIL"),
	}
}
