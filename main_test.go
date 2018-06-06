package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	app = App{}
	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME_TEST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
	)

	code := m.Run()

	droptable()

	os.Exit(code)
}
