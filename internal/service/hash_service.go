package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"

	"database/sql"
)

type HashService struct {
	db *sql.DB
	mu sync.Mutex
	N  int
}

func NewHashService(db *sql.DB, maxRequests int) *HashService {
	return &HashService{db: db, N: maxRequests}
}

func InitializeDatabase(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS hashes (
		hash TEXT PRIMARY KEY,
		description TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS request_limits (
		user_id INTEGER NOT NULL,
		request_count INTEGER NOT NULL,
		last_request_time TIMESTAMP NOT NULL,
		PRIMARY KEY (user_id)
	);
	`
	if _, err := db.Exec(query); err != nil {
		panic(fmt.Sprintf("Error initializing database: %v", err))
	}
}

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

	if !isEnglish(input) {
		return "", fmt.Errorf("input must contain only English letters")
	}

	hash := s.GenerateMD5(input)

	if _, err := s.db.Exec("INSERT INTO hashes (hash, description) VALUES ($1, $2)", hash, input); err != nil {
		return "", fmt.Errorf("error adding hash to database: %v", err)
	}

	return hash, nil
}

func (s *HashService) FindValueByHash(hash string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var description string
	err := s.db.QueryRow("SELECT description FROM hashes WHERE hash = $1", hash).Scan(&description)
	if err != nil {
		return "", fmt.Errorf("error finding value by hash in database: %v", err)
	}

	return description, nil
}

func (s *HashService) IncrementRequestCount(userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	var lastRequestTime string

	err := s.db.QueryRow("SELECT request_count, last_request_time FROM request_limits WHERE user_id = $1", userID).Scan(&count, &lastRequestTime)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		count = 0
		lastRequestTime = "1970-01-01 00:00:00"
	}

	now := time.Now()
	lastTime, _ := time.Parse("2006-01-02 15:04:05", lastRequestTime)

	if now.Sub(lastTime) > time.Hour {
		count = 0
	}

	count++

	if count > s.N {
		return fmt.Errorf("request limit exceeded")
	}

	if _, err := s.db.Exec("INSERT INTO request_limits (user_id, request_count, last_request_time) VALUES ($1, $2, $3) ON CONFLICT(user_id) DO UPDATE SET request_count = $2, last_request_time = $3", userID, count, now.Format("2006-01-02 15:04:05")); err != nil {
		return err
	}

	return nil
}

func (s *HashService) GenerateMD5(input string) string {
	hash := md5.Sum([]byte(strings.ToLower(input)))
	return hex.EncodeToString(hash[:])
}
