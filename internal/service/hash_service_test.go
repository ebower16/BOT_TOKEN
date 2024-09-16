package service

import (
	"strconv"
	"testing"
	"time"
)

func TestAddAndFindHash(t *testing.T) {
	hashService := NewHashService()

	inputText := "test string"
	hashValue, _ := hashService.AddHash(inputText)

	foundValue, err := hashService.FindValueByHash(hashValue)
	if err != nil || foundValue != inputText {
		t.Errorf("Expected '%s', got '%s', error: %v", inputText, foundValue, err)
	}
}

func TestGenerateMD5Performance(t *testing.T) {
	hashService := NewHashService()

	inputs := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		inputs[i] = "test string " + strconv.Itoa(i)
	}

	startTimeSingleThreaded := time.Now()
	for _, input := range inputs {
		hashService.GenerateMD5(input)
	}
	durationSingleThreaded := time.Since(startTimeSingleThreaded)

	startTimeMultiThreaded := time.Now()
	hashService.GenerateMD5Parallel(inputs)
	durationMultiThreaded := time.Since(startTimeMultiThreaded)

	t.Logf("Однопоточный вариант занял: %v", durationSingleThreaded)
	t.Logf("Многопоточный вариант занял: %v", durationMultiThreaded)

	if durationMultiThreaded >= durationSingleThreaded {
		t.Errorf("Параллельный вариант не быстрее однопоточного. Однопоточный: %v, Параллельный: %v", durationSingleThreaded, durationMultiThreaded)
	}
}
