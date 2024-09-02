package main

import (
	"testing"
)

func TestProcessMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Test secret", "secret", "md5(\"secret\") = 5ebe2294ecd0e0f179a11b0f600bfb0c\nreverse(\"5ebe2294ecd0e0f179a11b0f600bfb0c\") = \"secret\""},
		{"Test md5", "md5", "5d41402abc4b2a76b9719d911017c592"},
		{"Test invalid", "invalid", "Неверный текст. Пожалуйста, выберите опцию из кнопок или отправьте 'secret' или 'md5'."},
		{"Test empty", "", "Неверный текст. Пожалуйста, выберите опцию из кнопок или отправьте 'secret' или 'md5'."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := processMessage(test.input)
			if result != test.expected {
				t.Errorf("processMessage(%s) = %s, expected %s", test.input, result, test.expected)
			}
		})
	}
}

func TestGetKeyboard(t *testing.T) {
	keyboard := getKeyboard()

	if len(keyboard.Keyboard) != 3 {
		t.Errorf("getKeyboard() returned keyboard with %d rows, expected 3", len(keyboard.Keyboard))
	}

	if len(keyboard.Keyboard[0]) != 1 {
		t.Errorf("getKeyboard() returned keyboard with %d buttons in first row, expected 1", len(keyboard.Keyboard[0]))
	}

	if len(keyboard.Keyboard[1]) != 2 {
		t.Errorf("getKeyboard() returned keyboard with %d buttons in second row, expected 2", len(keyboard.Keyboard[1]))
	}

	if len(keyboard.Keyboard[2]) != 2 {
		t.Errorf("getKeyboard() returned keyboard with %d buttons in third row, expected 2", len(keyboard.Keyboard[2]))
	}
}
