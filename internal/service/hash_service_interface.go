package service

// HashServiceInterface defines methods for managing hashes.
type HashServiceInterface interface {
	AddHash(input string) (string, error)
	FindValueByHash(hash string) (string, error)
	IncrementRequestCount(userID int64) error
}
