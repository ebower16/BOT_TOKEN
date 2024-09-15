package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type HashService struct {
	hashes map[string]string
}

func NewHashService() *HashService {
	return &HashService{
		hashes: make(map[string]string),
	}
}

func (s *HashService) AddHash(input, hash string) {
	s.hashes[hash] = input
}

func (s *HashService) FindValueByHash(hash string) (string, error) {
	if value, exists := s.hashes[hash]; exists {
		return value, nil
	}
	return "", fmt.Errorf("значение не найдено для данного хеша")
}

func (s *HashService) GenerateMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
