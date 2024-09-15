package service

import (
	"testing"
)

func TestGenerateMD5(t *testing.T) {
	hashService := NewHashService()
	input := "test"
	expectedHash := "d8578edf8458ce06fbc5bb76a58c5ca" // MD5 для "test"

	hash := hashService.GenerateMD5(input)
	if hash != expectedHash {
		t.Errorf("Expected %s but got %s", expectedHash, hash)
	}
}

func TestAddHashAndFindValueByHash(t *testing.T) {
	hashService := NewHashService()
	input := "test"
	hash := hashService.GenerateMD5(input)
	hashService.AddHash(input, hash)

	value, err := hashService.FindValueByHash(hash)
	if err != nil {
		t.Fatalf("Expected no error but got %v", err)
	}

	if value != input {
		t.Errorf("Expected %s but got %s", input, value)
	}
}

func TestFindValueByHashNotFound(t *testing.T) {
	hashService := NewHashService()
	_, err := hashService.FindValueByHash("nonexistenthash")
	if err == nil {
		t.Fatal("Expected error but got none")
	}
}
