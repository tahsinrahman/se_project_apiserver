package main

import (
	"os"
)

var SecretKey string
var app App

func main() {
	SecretKey = os.Getenv("SECRET")

	app = App{}
	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
	)

	app.Run(":8080")
}
