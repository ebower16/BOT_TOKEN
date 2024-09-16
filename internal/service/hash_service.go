package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"database/sql"
)

type HashService struct {
	db *sql.DB
	mu sync.Mutex
}

func NewHashService(db *sql.DB) *HashService {
	return &HashService{
		db: db,
	}
}

// Проверка, что строка состоит только из английских букв.
func isEnglish(input string) bool {
	for _, char := range input {
		if !(char >= 'a' && char <= 'z') && !(char >= 'A' && char <= 'Z') {
			return false
		}
	}
	return true
}

func (s *HashService) AddHash(input string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !isEnglish(input) { // Проверяем, что входные данные состоят только из английских букв.
		return "", fmt.Errorf("входные данные должны содержать только английские буквы")
	}

	hash := s.GenerateMD5(input)

	if _, err := s.db.Exec("INSERT INTO hashes (hash, description) VALUES ($1, $2)", hash, input); err != nil {
		return "", fmt.Errorf("ошибка добавления хеша в базу данных: %v", err)
	}

	return hash, nil
}

func (s *HashService) FindValueByHash(hash string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var description string
	err := s.db.QueryRow("SELECT description FROM hashes WHERE hash = $1", hash).Scan(&description)
	if err != nil {
		return "", fmt.Errorf("ошибка поиска значения по хешу в базе данных: %v", err)
	}

	return description, nil
}

func (s *HashService) GenerateMD5(input string) string {
	hash := md5.Sum([]byte(strings.ToLower(input)))
	return hex.EncodeToString(hash[:])
}
