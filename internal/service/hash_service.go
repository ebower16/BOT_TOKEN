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

// AddHash добавляет хеш в память.
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

// FindValueByHash ищет значение по хешу.
func (s *HashService) FindValueByHash(hash string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if value, exists := s.hashes[hash]; exists {
		return value.Description, nil
	}
	return "", fmt.Errorf("значение не найдено для данного хеша")
}

// GenerateMD5 генерирует MD5-хеш для строки.
func (s *HashService) GenerateMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// GenerateMD5Parallel генерирует MD5-хеши для нескольких строк параллельно.
func (s *HashService) GenerateMD5Parallel(inputs []string) []string {
	var wg sync.WaitGroup
	hashes := make([]string, len(inputs))
	mu := sync.Mutex{}

	for i, input := range inputs {
		wg.Add(1)
		go func(i int, input string) {
			defer wg.Done()
			hash := s.GenerateMD5(input)
			mu.Lock()
			hashes[i] = hash
			mu.Unlock()
		}(i, input)
	}

	wg.Wait()
	return hashes
}
