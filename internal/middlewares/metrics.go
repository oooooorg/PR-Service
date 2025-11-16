package middlewares

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{0.001, 0.01, 0.05, 0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "Total number of HTTP request errors",
		},
		[]string{"method", "path", "status"},
	)

	httpActiveRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_active_requests",
			Help: "Number of active HTTP requests being processed",
		},
	)

	dbConnectionsInUse = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_in_use",
			Help: "Number of database connections currently in use",
		},
	)

	dbConnectionsOpen = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_open",
			Help: "Number of open database connections",
		},
	)
)

func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			httpActiveRequests.Inc()
			defer httpActiveRequests.Dec()

			err := next(c)

			status := c.Response().Status
			if err != nil {
				var he *echo.HTTPError
				if errors.As(err, &he) {
					status = he.Code
				}
			}

			method := c.Request().Method
			path := c.Path()
			statusStr := strconv.Itoa(status)

			duration := time.Since(start).Seconds()
			httpRequestDuration.WithLabelValues(method, path, statusStr).Observe(duration)

			httpRequestsTotal.WithLabelValues(method, path, statusStr).Inc()

			if status >= 400 {
				httpRequestErrors.WithLabelValues(method, path, statusStr).Inc()
			}

			return err
		}
	}
}

func UpdateDBMetrics(db *sql.DB) {
	stats := db.Stats()
	dbConnectionsInUse.Set(float64(stats.InUse))
	dbConnectionsOpen.Set(float64(stats.OpenConnections))
}
