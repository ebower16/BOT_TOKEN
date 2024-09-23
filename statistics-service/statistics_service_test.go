package statistics

import (
	"testing"
)

func TestStatisticsIncrement(t *testing.T) {
	stats := NewStatistics()

	stats.Increment()

	if count := stats.GetCount(); count != 1 {
		t.Fatalf("Expected count to be 1, got %d", count)
	}
}
