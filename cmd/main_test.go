package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessMessage(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Secret",
			input:    "secret",
			expected: "md5(\"secret\") = 5ebe2294ecd0e0f08eab7690d2a6ee69\nreverse(\"5ebe2294ecd0e0f08eab7690d2a6ee69\") = \"secret\"",
		},
		{
			name:     "MD5",
			input:    "md5",
			expected: "5d41402abc4b2a76b9719d911017c592",
		},
		{
			name:     "Invalid",
			input:    "invalid",
			expected: "Неверный текст. Пожалуйста, выберите опцию из кнопок или отправьте 'secret' или 'md5'.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := processMessage(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
