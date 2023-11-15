package config

import (
	"github.com/joho/godotenv"
	"os"
)

var SECRET_KEY string

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	SECRET_KEY = os.Getenv("SECRET_KEY")
}
