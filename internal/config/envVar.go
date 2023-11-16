package config

import (
	"github.com/joho/godotenv"
	"os"
)

var SecretKey string

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	SecretKey = os.Getenv("SECRET_KEY")
}
