package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CacheOperations tracks the number of cache operations by type and status
	CacheOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_operations_total",
			Help: "The total number of cache operations",
		},
		[]string{"operation", "status"},
	)

	// CacheSize tracks the current number of items in cache
	CacheSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cache_items_current",
			Help: "The current number of items in cache",
		},
	)

	// CacheHits tracks cache hits and misses
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "The total number of cache hits/misses",
		},
		[]string{"result"},
	)

	// RequestDuration tracks the duration of HTTP requests
	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)
