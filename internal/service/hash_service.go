package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"botus/internal/domain"
	"database/sql"
)

type HashService struct {
	db     *sql.DB
	mu     sync.Mutex
	hashes map[string]domain.Hash
}

func NewHashService(db *sql.DB) *HashService {
	return &HashService{
		db:     db,
		hashes: make(map[string]domain.Hash),
	}
}

func (s *HashService) AddHash(input string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash := s.GenerateMD5(input)
	s.hashes[hash] = domain.Hash{
		ID:          len(s.hashes) + 1, // Пример ID, можно использовать автоинкремент в базе данных
		Time:        time.Now(),
		Description: input,
	}

	return hash
}

func (s *HashService) FindValueByHash(hash string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if value, exists := s.hashes[hash]; exists {
		return value.Description, nil
	}
	return "", fmt.Errorf("значение не найдено для данного хеша")
}

func (s *HashService) GenerateMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
