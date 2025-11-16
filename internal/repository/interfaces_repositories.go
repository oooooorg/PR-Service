package repository

import (
	"context"
	"database/sql"

	"github.com/oooooorg/PR-Service/internal/entity"
)

type UserRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateUser(ctx context.Context, tx *sql.Tx, user *entity.User) error
	GetUsersByTeam(ctx context.Context, tx *sql.Tx, teamName string) ([]*entity.User, error)
	GetUserByID(ctx context.Context, tx *sql.Tx, userID string) (*entity.User, error)
	SetUserActive(ctx context.Context, tx *sql.Tx, userID string, isActive bool) (*entity.User, error)
}

type TeamRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateTeam(ctx context.Context, tx *sql.Tx, team *entity.Team) error
	GetTeamByName(ctx context.Context, tx *sql.Tx, teamName string) (*entity.Team, error)
	TeamExists(ctx context.Context, tx *sql.Tx, teamName string) (bool, error)
}

type PullRequestRepository interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreatePullRequest(ctx context.Context, tx *sql.Tx, pr *entity.PullRequest) error
	GetPullRequestByID(ctx context.Context, tx *sql.Tx, prID string) (*entity.PullRequest, error)
	UpdatePullRequestStatus(ctx context.Context, tx *sql.Tx, prID string, status string) (*entity.PullRequest, error)
	UpdatePullRequestReviewers(ctx context.Context, tx *sql.Tx, prID string, reviewer1, reviewer2 string) (*entity.PullRequest, error)
	GetPullRequestsByReviewer(ctx context.Context, tx *sql.Tx, reviewerID string) ([]*entity.PullRequest, error)
}
