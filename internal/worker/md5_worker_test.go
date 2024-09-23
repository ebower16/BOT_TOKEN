package worker

import (
	"crypto/md5"
	"encoding/hex"
	"testing"
	"time"
)

func TestComputeMD5Performance(t *testing.T) {
	values := []string{"test1", "test2", "test3", "test4", "test5"}

	
	start := time.Now()
	singleResults := make(map[string]string)

	for _, value := range values {
		hash := md5.Sum([]byte(value))
		singleResults[value] = hex.EncodeToString(hash[:])
	}

	elapsedSingle := time.Since(start)

	
	start = time.Now()
	parallelWorker := NewHashWorker()
	parallelWorker.ComputeMD5(values)

	elapsedParallel := time.Since(start)

	t.Logf("Single-threaded duration: %v", elapsedSingle)
	t.Logf("Multi-threaded duration: %v", elapsedParallel)

	if elapsedParallel >= elapsedSingle {
		t.Error("Parallel computation is not faster than single-threaded")
	}
}
