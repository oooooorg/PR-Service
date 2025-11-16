package middlewares

import (
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

func LoggerMiddleware(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			latency := time.Since(start)
			req := c.Request()
			res := c.Response()

			logger.Info("HTTP request",
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.Int("status", res.Status),
				slog.Duration("latency", latency),
			)

			return err
		}
	}
}
