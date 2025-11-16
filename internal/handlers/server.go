package handlers

import (
	"database/sql"
	"log/slog"

	"github.com/oooooorg/PR-Service/internal/config"
	"github.com/oooooorg/PR-Service/internal/repository"
	"github.com/oooooorg/PR-Service/internal/service"
)

type Server struct {
	db                 *sql.DB
	cfg                *config.Config
	pullRequestService service.PullRequestService
	teamService        service.TeamService
	userService        service.UserService
	logger             *slog.Logger
}

func NewServer(logger *slog.Logger, db *sql.DB, cfg *config.Config) *Server {
	userRepository := repository.NewUserRepository(db)
	teamRepository := repository.NewTeamRepository(db)
	pullRequestRepository := repository.NewPullRequestRepository(db)

	return &Server{
		logger:             logger,
		db:                 db,
		cfg:                cfg,
		pullRequestService: service.NewPullRequestService(logger, pullRequestRepository, userRepository, teamRepository),
		teamService:        service.NewTeamService(logger, userRepository, teamRepository),
		userService:        service.NewUserService(logger, userRepository, teamRepository),
	}
}
