package middleware

import (
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// Basit performans metrikleri (in-memory)
var (
	totalRequests  uint64
	total5xx       uint64
	total4xx       uint64
	total2xx       uint64
	totalLatencyNs uint64
)

// PerformanceMonitor metrik toplar.
func PerformanceMonitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start)
		atomic.AddUint64(&totalRequests, 1)
		atomic.AddUint64(&totalLatencyNs, uint64(dur.Nanoseconds()))
		status := c.Writer.Status()
		switch {
		case status >= 500:
			atomic.AddUint64(&total5xx, 1)
		case status >= 400:
			atomic.AddUint64(&total4xx, 1)
		default:
			atomic.AddUint64(&total2xx, 1)
		}
	}
}

// PerfStats anlık metrik snapshot’ı
type PerfStats struct {
	TotalRequests uint64  `json:"total_requests"`
	Total2xx      uint64  `json:"total_2xx"`
	Total4xx      uint64  `json:"total_4xx"`
	Total5xx      uint64  `json:"total_5xx"`
	AvgLatencyMs  float64 `json:"avg_latency_ms"`
}

// GetPerfStats döner
func GetPerfStats() PerfStats {
	tr := atomic.LoadUint64(&totalRequests)
	t2 := atomic.LoadUint64(&total2xx)
	t4 := atomic.LoadUint64(&total4xx)
	t5 := atomic.LoadUint64(&total5xx)
	lat := atomic.LoadUint64(&totalLatencyNs)
	var avg float64
	if tr > 0 {
		avg = float64(lat) / float64(tr) / 1e6
	}
	return PerfStats{TotalRequests: tr, Total2xx: t2, Total4xx: t4, Total5xx: t5, AvgLatencyMs: avg}
}
