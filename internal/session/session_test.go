package session

import (
	"testing"
	"time"
)

func TestUserSession(t *testing.T) {
	session := NewUserSession(1, "testuser")

	if session.UserID != 1 {
		t.Errorf("Expected UserID 1, got %d", session.UserID)
	}

	if session.Username != "testuser" {
		t.Errorf("Expected Username 'testuser', got '%s'", session.Username)
	}

	session.IncrementAttempts()
	if session.Attempts != 1 {
		t.Errorf("Expected Attempts 1, got %d", session.Attempts)
	}

	session.ResetAttempts()
	if session.Attempts != 0 {
		t.Errorf("Expected Attempts 0 after reset, got %d", session.Attempts)
	}

	session.BlockUser(10 * time.Second)
	if !session.IsBlocked() {
		t.Error("Expected session to be blocked")
	}

	time.Sleep(11 * time.Second)
	if session.IsBlocked() {
		t.Error("Expected session to be unblocked after
