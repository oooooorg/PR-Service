package service

import (
	"context"
	"errors"
	"log/slog"

	api "github.com/oooooorg/PR-Service/internal/gen"
	"github.com/oooooorg/PR-Service/internal/models"
	"github.com/oooooorg/PR-Service/internal/repository"
)

var ErrUserNotFound = errors.New("team already exists")
var ErrUserExists = errors.New("user already exists")

type UserServiceImpl struct {
	logger   *slog.Logger
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewUserService(
	logger *slog.Logger,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
) UserService {
	return &UserServiceImpl{
		logger:   logger,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (u *UserServiceImpl) SetUserActive(ctx context.Context, req *api.PostUsersSetIsActiveJSONRequestBody) (*models.User, error) {
	if req.UserId == "" {
		return nil, errors.New("UserId is required")
	}

	tx, err := u.userRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = u.userRepo.GetUserByID(ctx, tx, req.UserId)
	if err != nil {
		return nil, err
	}

	updatedUser, err := u.userRepo.SetUserActive(ctx, tx, req.UserId, req.IsActive)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	resultUser := &models.User{
		UserId:   updatedUser.UserID,
		Username: updatedUser.Username,
		TeamName: updatedUser.TeamName,
		IsActive: updatedUser.IsActive,
	}

	return resultUser, nil
}
