package database

import (
	"testing"

	_ "github.com/mattn/go-sqlite3" // Импортируем драйвер SQLite
)

func TestInitDatabase(t *testing.T) {
	db := InitDatabase()
	if db == nil {
		t.Fatal("Expected database to be initialized, got nil")
	}
	defer db.Close()
}

func TestRegisterUser(t *testing.T) {
	db := InitDatabase()
	defer db.Close()

	username := "testuser"
	password := "testpassword"
	chatID := int64(123456)

	if !RegisterUser(username, password, db, chatID) {
		t.Fatal("Expected successful registration, got failure")
	}

	if RegisterUser(username, password, db, chatID) {
		t.Fatal("Expected registration failure for existing user, got success")
	}
}

func TestCheckUserCredentials(t *testing.T) {
	db := InitDatabase()
	defer db.Close()

	username := "testuser"
	password := "testpassword"
	chatID := int64(123456)

	RegisterUser(username, password, db, chatID)

	if !CheckUserCredentials(username, password, db) {
		t.Fatal("Expected successful login, got failure")
	}

	if CheckUserCredentials(username, "wrongpassword", db) {
		t.Fatal("Expected login failure for wrong password, got success")
	}
}

func TestGetAllUsers(t *testing.T) {
	db := InitDatabase()
	defer db.Close()

	RegisterUser("user1", "password1", db, 123456)
	RegisterUser("user2", "password2", db, 789012)

	users, err := GetAllUsers(db)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}
