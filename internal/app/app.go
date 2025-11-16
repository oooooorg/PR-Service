package app

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/oooooorg/PR-Service/internal/config"
	api "github.com/oooooorg/PR-Service/internal/gen"
	"github.com/oooooorg/PR-Service/internal/handlers"
	"github.com/oooooorg/PR-Service/internal/middlewares"
)

type App struct {
	db     *sql.DB
	cfg    *config.Config
	logger *slog.Logger
}

func New(cfg *config.Config, db *sql.DB, logger *slog.Logger) *App {
	return &App{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

func (app *App) Run() error {
	echoApp := echo.New()
	echoApp.HideBanner = true

	echoApp.Use(
		middleware.Recover(),
		middleware.CORS(),
		middlewares.LoggerMiddleware(app.logger),
		middlewares.PrometheusMiddleware(),
	)

	server := handlers.NewServer(app.logger, app.db, app.cfg)
	api.RegisterHandlers(echoApp, server)

	echoApp.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	routes := echoApp.Routes()
	for _, route := range routes {
		app.logger.Info("Registered route",
			slog.String("method", route.Method),
			slog.String("path", route.Path),
		)
	}

	stopMetrics := make(chan struct{})
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				middlewares.UpdateDBMetrics(app.db)
			case <-stopMetrics:
				app.logger.Info("Stopping metrics updater")
				return
			}
		}
	}()

	go func() {
		app.logger.Info("Starting HTTP server", slog.String("port", "8080"))

		if err := echoApp.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("Server error", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	app.logger.Info("Received shutdown signal", slog.String("signal", sig.String()))

	close(stopMetrics)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := echoApp.Shutdown(ctx); err != nil {
		app.logger.Error("Server shutdown error", slog.String("error", err.Error()))
		return err
	}

	if err := app.db.Close(); err != nil {
		app.logger.Error("Database close error", slog.String("error", err.Error()))
		return err
	}

	app.logger.Info("Server stopped gracefully")
	return nil
}
