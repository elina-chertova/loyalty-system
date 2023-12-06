package config

import (
	"os"

	"github.com/joho/godotenv"
)

var SecretKey string

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	SecretKey = os.Getenv("SECRET_KEY")
}
