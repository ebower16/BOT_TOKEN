package service

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateMD5Performance(t *testing.T) {
	hashService := NewHashService(nil) // Передайте nil, если база данных не используется в тесте

	// Создаем массив строк для тестирования
	inputs := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		inputs[i] = fmt.Sprintf("test string %d", i)
	}

	// Однопоточный тест
	start := time.Now()
	for _, input := range inputs {
		hashService.GenerateMD5(input)
	}
	oneThreadDuration := time.Since(start)

	// Параллельный тест
	start = time.Now()
	hashService.GenerateMD5Parallel(inputs)
	parallelDuration := time.Since(start)

	// Выводим результаты
	t.Logf("Однопоточный вариант занял: %v", oneThreadDuration)
	t.Logf("Параллельный вариант занял: %v", parallelDuration)

	// Проверяем, что параллельный вариант быстрее
	if parallelDuration >= oneThreadDuration {
		t.Errorf("Параллельный вариант не быстрее однопоточного. Однопоточный: %v, Параллельный: %v", oneThreadDuration, parallelDuration)
	}
}
