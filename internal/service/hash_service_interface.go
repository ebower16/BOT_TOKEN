package service


type HashServiceInterface interface {
	AddHash(input string) (string, error)
	FindValueByHash(hash string) (string, error)
	IncrementRequestCount(userID int64) error
}
