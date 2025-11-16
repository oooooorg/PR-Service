package service

import (
	"context"

	api "github.com/oooooorg/PR-Service/internal/gen"
	"github.com/oooooorg/PR-Service/internal/models"
)

type PullRequestService interface {
	CreatePullRequest(ctx context.Context, req *api.PostPullRequestCreateJSONRequestBody) (*models.PullRequest, error)
	MergePullRequest(ctx context.Context, req *api.PostPullRequestMergeJSONRequestBody) (*models.PullRequest, error)
	ReassignReviewer(ctx context.Context, req *api.PostPullRequestReassignJSONRequestBody) (*models.PullRequest, string, error)
	GetUserReviewRequests(ctx context.Context, req *api.GetUsersGetReviewParams) ([]*models.PullRequestShort, error)
}

type TeamService interface {
	CreateTeam(ctx context.Context, team *api.Team) (*models.Team, error)
	GetTeam(ctx context.Context, req *api.GetTeamGetParams) (*models.Team, error)
}

type UserService interface {
	SetUserActive(ctx context.Context, req *api.PostUsersSetIsActiveJSONRequestBody) (*models.User, error)
}
