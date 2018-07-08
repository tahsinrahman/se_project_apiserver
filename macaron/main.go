package main

import (
	"os"
)

var SecretKey string

func main() {
	SecretKey = os.Getenv("SECRET")

	InitializeDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
	)

	InitializeRouter()
}
