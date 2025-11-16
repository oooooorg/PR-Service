package main

import (
	"log/slog"
	"os"

	_ "github.com/lib/pq"

	"github.com/oooooorg/PR-Service/internal/app"
	"github.com/oooooorg/PR-Service/internal/config"
	"github.com/oooooorg/PR-Service/internal/database"
	"github.com/oooooorg/PR-Service/internal/logger"
)

const ConfigPath = "config.yml"

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	log := logger.New(env)
	slog.SetDefault(log)

	log.Info("Starting PR-Service",
		slog.String("env", env),
		slog.String("version", "1.0.0"),
	)

	cfg, err := config.NewConfig(ConfigPath)
	if err != nil {
		log.Error("Failed to load config",
			slog.String("error", err.Error()),
			slog.String("path", ConfigPath),
		)
		os.Exit(1)
	}

	log.Info("Config loaded successfully")

	log.Info("Connecting to database",
		slog.String("host", cfg.Database.Host),
		slog.String("port", cfg.Database.Port),
		slog.String("database", cfg.Database.DBName),
		slog.String("user", cfg.Database.User),
	)

	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Error("Failed to connect to database",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	log.Info("Database connection established")

	application := app.New(cfg, db, log)

	if err := application.Run(); err != nil {
		log.Error("Application error",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	log.Info("Application stopped successfully")
}
