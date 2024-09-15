package domain

import (
	"crypto/md5"
	"encoding/hex"
)

// Hash представляет MD5 хеш
type Hash struct {
	Value string
}

// NewHash создает новый хеш из строки
func NewHash(input string) *Hash {
	hash := md5.Sum([]byte(input))
	return &Hash{Value: hex.EncodeToString(hash[:])}
}

// Equals проверяет, равен ли текущий хеш другому хешу
func (h *Hash) Equals(other *Hash) bool {
	return h.Value == other.Value
}

// String возвращает строковое представление хеша
func (h *Hash) String() string {
	return h.Value
}
