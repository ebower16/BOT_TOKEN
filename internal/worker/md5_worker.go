package worker

import (
	"crypto/md5"
	"encoding/hex"
	"sync"
)

type HashWorker struct {
	results map[string]string
	mu      sync.Mutex
}

func NewHashWorker() *HashWorker {
	return &HashWorker{
		results: make(map[string]string),
	}
}

func (hw *HashWorker) ComputeMD5(values []string) {
	var wg sync.WaitGroup

	for _, value := range values {
		wg.Add(1)
		go func(val string) {
			defer wg.Done()
			hash := md5.Sum([]byte(val))
			hw.mu.Lock()
			hw.results[val] = hex.EncodeToString(hash[:])
			hw.mu.Unlock()
		}(value)
	}

	wg.Wait()
}

func (hw *HashWorker) GetResults() map[string]string {
	return hw.results
}
