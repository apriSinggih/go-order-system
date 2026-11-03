package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

// InitMetrics memasang middleware HTTP metrics dan mengekspos endpoint /metrics
// dalam format Prometheus. Path /metrics dan /health dikecualikan dari
// perhitungan agar tidak mendistorsi data RPS/latency.
func InitMetrics(r *gin.Engine) {
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10})
	m.SetExcludePaths([]string{"/metrics", "/health"})
	m.Use(r)
}
