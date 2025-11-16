package handlers

import (
	"github.com/labstack/echo/v4"

	api "github.com/oooooorg/PR-Service/internal/gen"
)

type PullRequestService interface {
	CreatePullRequest(ctx echo.Context) error
	MergePullRequest(ctx echo.Context) error
	ReassignPullRequestReviewer(ctx echo.Context) error
	GetUserReviewRequests(ctx echo.Context, params api.GetUsersGetReviewParams) error
}

type TeamService interface {
	CreateTeam(ctx echo.Context) error
	GetTeam(ctx echo.Context, params api.GetTeamGetParams) error
}

type UserService interface {
	SetUserActive(ctx echo.Context) error
}

type ServiceHandler interface {
	PullRequestService
	TeamService
	UserService
}
