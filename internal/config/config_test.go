package config

import (
	"os"
	"testing"
)

func TestLoadSuccess(t *testing.T) {

	err := os.WriteFile(".env", []byte("DATABASE_URL=postgres://user:pass@localhost/db\nAPI_KEY=123456"), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")

	err = Load()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if os.Getenv("DATABASE_URL") != "postgres://user:pass@localhost/db" {
		t.Errorf("Expected DATABASE_URL to be 'postgres://user:pass@localhost/db', got '%s'", os.Getenv("DATABASE_URL"))
	}
	if os.Getenv("API_KEY") != "123456" {
		t.Errorf("Expected API_KEY to be '123456', got '%s'", os.Getenv("API_KEY"))
	}
}

func TestLoadMissingEnvFile(t *testing.T) {

	os.Remove(".env")

	err := Load()
	if err == nil {
		t.Fatal("Expected an error when .env file is missing, got nil")
	}
}

func TestLoadMissingRequiredVar(t *testing.T) {

	err := os.WriteFile(".env", []byte("DATABASE_URL=postgres://user:pass@localhost/db\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer os.Remove(".env")
	err = Load()
	if err == nil {
		t.Fatal("Expected an error due to missing API_KEY, got nil")
	}

	expectedError := "required environment variable API_KEY is not set"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
	}
}
