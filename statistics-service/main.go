package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of requests received",
		},
		[]string{"method"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Duration of requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

type Stats struct {
	sync.Mutex
	Count     int       `json:"count"`
	StartTime time.Time `json:"start_time"`
}

var stats Stats

func init() {
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

func main() {
	stats.StartTime = time.Now()

	r := mux.NewRouter()

	r.HandleFunc("/stats", getStats).Methods("GET")

	r.Handle("/metrics", promhttp.Handler())

	r.Use(middleware)

	http.ListenAndServe(":8080", r)
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()

		requestCounter.WithLabelValues(r.Method).Inc()
		requestDuration.WithLabelValues(r.Method).Observe(duration)

		stats.Lock()
		stats.Count++
		stats.Unlock()
	})
}

func getStats(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"count":      stats.Count,
		"start_time": stats.StartTime,
		"uptime":     time.Since(stats.StartTime).Seconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
