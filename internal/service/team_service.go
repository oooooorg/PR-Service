package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/oooooorg/PR-Service/internal/entity"
	api "github.com/oooooorg/PR-Service/internal/gen"
	"github.com/oooooorg/PR-Service/internal/models"
	"github.com/oooooorg/PR-Service/internal/repository"
)

var ErrTeamExists = errors.New("team already exists")
var ErrTeamNotFound = errors.New("team not found")

type TeamServiceImpl struct {
	logger   *slog.Logger
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewTeamService(
	logger *slog.Logger,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
) TeamService {
	return &TeamServiceImpl{
		logger:   logger,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (t *TeamServiceImpl) CreateTeam(ctx context.Context, team *api.Team) (*models.Team, error) {
	if team.TeamName == "" {
		return nil, errors.New("team name is required")
	}
	if len(team.Members) == 0 {
		return nil, errors.New("team must have at least one member")
	}

	tx, err := t.userRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	exists, err := t.teamRepo.TeamExists(ctx, tx, team.TeamName)
	if err != nil {
		return nil, fmt.Errorf("failed to check team existence: %w", err)
	}
	if exists {
		return nil, ErrTeamExists
	}

	teamEntity := &entity.Team{
		TeamName: team.TeamName,
	}

	err = t.teamRepo.CreateTeam(ctx, tx, teamEntity)
	if err != nil {
		return nil, err
	}

	for _, member := range team.Members {
		if member.UserId == "" {
			return nil, errors.New("user ID is required for all members")
		}

		_, err := t.userRepo.GetUserByID(ctx, tx, member.UserId)
		if err == nil {
			return nil, err
		}

		userEntity := &entity.User{
			UserID:   member.UserId,
			Username: member.Username,
			IsActive: member.IsActive,
			TeamName: team.TeamName,
		}

		err = t.userRepo.CreateUser(ctx, tx, userEntity)
		if err != nil {
			return nil, ErrUserExists
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	createdTeam, err := t.GetTeam(ctx, &models.TeamGetParams{
		TeamName: team.TeamName,
	})

	if err != nil {
		return nil, err
	}

	return createdTeam, nil
}

func (t *TeamServiceImpl) GetTeam(ctx context.Context, req *api.GetTeamGetParams) (*models.Team, error) {
	if req.TeamName == "" {
		return nil, fmt.Errorf("team name is required")
	}

	tx, err := t.teamRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	teamEntity, err := t.teamRepo.GetTeamByName(ctx, tx, req.TeamName)
	if err != nil {
		return nil, err
	}

	users, err := t.userRepo.GetUsersByTeam(ctx, tx, teamEntity.TeamName)
	if err != nil {
		return nil, err
	}

	members := make([]models.TeamMember, 0, len(users))
	for _, u := range users {
		members = append(members, models.TeamMember{
			IsActive: u.IsActive,
			UserId:   u.UserID,
			Username: u.Username,
		})
	}

	team := &models.Team{
		TeamName: teamEntity.TeamName,
		Members:  members,
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return team, nil
}
