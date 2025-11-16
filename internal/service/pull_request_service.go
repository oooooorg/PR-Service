package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"github.com/oooooorg/PR-Service/internal/entity"
	api "github.com/oooooorg/PR-Service/internal/gen"
	"github.com/oooooorg/PR-Service/internal/models"
	"github.com/oooooorg/PR-Service/internal/repository"
)

var ErrPullRequestExists = errors.New("pull request already exists")
var ErrPullRequestNotAsigned = errors.New("pull request not_asigned")
var ErrPullRequestNoCandidate = errors.New("pull request no_candidate")
var ErrPullRequestMerged = errors.New("pull request already merged")

type PullRequestServiceImpl struct {
	logger   *slog.Logger
	prRepo   repository.PullRequestRepository
	userRepo repository.UserRepository
	teamRepo repository.TeamRepository
}

func NewPullRequestService(
	logger *slog.Logger,
	prRepo repository.PullRequestRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
) PullRequestService {
	return &PullRequestServiceImpl{
		logger:   logger,
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (p *PullRequestServiceImpl) GetUsersForPR(ctx context.Context, authorID string) ([]string, error) {
	if authorID == "" {
		return nil, errors.New("authorID is required")
	}

	tx, err := p.prRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	author, err := p.userRepo.GetUserByID(ctx, tx, authorID)
	if err != nil {
		return nil, err
	}

	users, err := p.userRepo.GetUsersByTeam(ctx, tx, author.TeamName)
	if err != nil {
		return nil, err
	}

	var candidates []*entity.User
	for _, u := range users {
		if !u.IsActive {
			continue
		}
		if u.UserID == authorID {
			continue
		}
		candidates = append(candidates, u)
	}

	if len(candidates) == 0 {
		return []string{}, nil
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	n := 2
	if len(candidates) < n {
		n = len(candidates)
	}

	reviewers := make([]string, n)
	for i := 0; i < n; i++ {
		reviewers[i] = candidates[i].UserID
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return reviewers, nil
}

func (p *PullRequestServiceImpl) CreatePullRequest(ctx context.Context, req *api.PostPullRequestCreateJSONRequestBody) (*models.PullRequest, error) {
	if req.AuthorId == "" {
		return nil, errors.New("AuthorId is empty")
	}
	if req.PullRequestId == "" {
		return nil, errors.New("PullRequestId is empty")
	}
	if req.PullRequestName == "" {
		return nil, errors.New("PullRequestName is empty")
	}

	tx, err := p.prRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = p.prRepo.GetPullRequestByID(ctx, tx, req.PullRequestId)
	if err == nil {
		return nil, ErrPullRequestExists
	}

	user, err := p.userRepo.GetUserByID(ctx, tx, req.AuthorId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	_, err = p.teamRepo.GetTeamByName(ctx, tx, user.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}

	reviewers, err := p.GetUsersForPR(ctx, req.AuthorId)
	if err != nil {
		return nil, err
	}

	var reviewer1, reviewer2 string
	if len(reviewers) > 0 {
		reviewer1 = reviewers[0]
	}
	if len(reviewers) > 1 {
		reviewer2 = reviewers[1]
	}

	pullRequestEntity := &entity.PullRequest{
		AuthorID:                req.AuthorId,
		PullRequestID:           req.PullRequestId,
		PullRequestName:         req.PullRequestName,
		Status:                  entity.StatusOpen,
		MergedAt:                nil,
		AssignedReviewersFirst:  reviewer1,
		AssignedReviewersSecond: reviewer2,
	}

	err = p.prRepo.CreatePullRequest(ctx, tx, pullRequestEntity)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	pullRequest := &models.PullRequest{
		PullRequestId:     pullRequestEntity.PullRequestID,
		PullRequestName:   pullRequestEntity.PullRequestName,
		AuthorId:          pullRequestEntity.AuthorID,
		Status:            api.PullRequestStatus(pullRequestEntity.Status),
		AssignedReviewers: []string{pullRequestEntity.AssignedReviewersFirst, pullRequestEntity.AssignedReviewersSecond},
		CreatedAt:         &pullRequestEntity.CreatedAt,
	}

	return pullRequest, nil
}

func (p *PullRequestServiceImpl) MergePullRequest(ctx context.Context, req *api.PostPullRequestMergeJSONRequestBody) (*models.PullRequest, error) {
	if req.PullRequestId == "" {
		return nil, errors.New("PullRequestId is empty")
	}

	tx, err := p.prRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	pr, err := p.prRepo.GetPullRequestByID(ctx, tx, req.PullRequestId)
	if err != nil {
		return nil, err
	}

	if pr == nil {
		return nil, sql.ErrNoRows
	}

	updatedPR, err := p.prRepo.UpdatePullRequestStatus(ctx, tx, req.PullRequestId, "MERGED")
	if err != nil {
		return nil, err
	}

	pullRequest := &models.PullRequest{
		PullRequestId:     updatedPR.PullRequestID,
		PullRequestName:   updatedPR.PullRequestName,
		AuthorId:          updatedPR.AuthorID,
		Status:            api.PullRequestStatus(updatedPR.Status),
		AssignedReviewers: []string{updatedPR.AssignedReviewersFirst, updatedPR.AssignedReviewersSecond},
		CreatedAt:         &updatedPR.CreatedAt,
		MergedAt:          updatedPR.MergedAt,
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pullRequest, nil
}

func (p *PullRequestServiceImpl) ReassignReviewer(ctx context.Context, req *api.PostPullRequestReassignJSONRequestBody) (*models.PullRequest, string, error) {
	if req.PullRequestId == "" {
		return nil, "", errors.New("PullRequestId is empty")
	}
	if req.OldUserId == "" {
		return nil, "", errors.New("OldUserId is empty")
	}

	tx, err := p.prRepo.BeginTx(ctx)
	if err != nil {
		return nil, "", err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	pr, err := p.prRepo.GetPullRequestByID(ctx, tx, req.PullRequestId)
	if err != nil {
		return nil, "", err
	}

	if pr.Status == entity.StatusMerged {
		return nil, "", ErrPullRequestMerged
	}

	var otherReviewer string

	if pr.AssignedReviewersFirst == req.OldUserId {
		otherReviewer = pr.AssignedReviewersSecond
	} else if pr.AssignedReviewersSecond == req.OldUserId {
		otherReviewer = pr.AssignedReviewersFirst
	} else {
		return nil, "", ErrPullRequestNotAsigned
	}

	author, err := p.userRepo.GetUserByID(ctx, tx, pr.AuthorID)
	if err != nil {
		return nil, "", err
	}

	users, err := p.userRepo.GetUsersByTeam(ctx, tx, author.TeamName)
	if err != nil {
		return nil, "", err
	}

	var candidates []*entity.User
	for _, u := range users {
		if !u.IsActive {
			continue
		}
		if u.UserID == pr.AuthorID {
			continue
		}
		if u.UserID == req.OldUserId {
			continue
		}
		if u.UserID == otherReviewer {
			continue
		}
		candidates = append(candidates, u)
	}

	if len(candidates) == 0 {
		return nil, "", ErrPullRequestNoCandidate
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	newReviewer := candidates[r.Intn(len(candidates))].UserID

	var updatedPR *entity.PullRequest
	if pr.AssignedReviewersFirst == req.OldUserId {
		updatedPR, err = p.prRepo.UpdatePullRequestReviewers(ctx, tx, req.PullRequestId, newReviewer, otherReviewer)
	} else {
		updatedPR, err = p.prRepo.UpdatePullRequestReviewers(ctx, tx, req.PullRequestId, otherReviewer, newReviewer)
	}

	if err != nil {
		return nil, "", err
	}

	pullRequest := &models.PullRequest{
		PullRequestId:     updatedPR.PullRequestID,
		PullRequestName:   updatedPR.PullRequestName,
		AuthorId:          updatedPR.AuthorID,
		Status:            api.PullRequestStatus(updatedPR.Status),
		AssignedReviewers: []string{updatedPR.AssignedReviewersFirst, updatedPR.AssignedReviewersSecond},
	}

	if err := tx.Commit(); err != nil {
		return nil, "", err
	}

	return pullRequest, newReviewer, nil
}

func (p *PullRequestServiceImpl) GetUserReviewRequests(ctx context.Context, req *api.GetUsersGetReviewParams) ([]*models.PullRequestShort, error) {
	if req.UserId == "" {
		return nil, errors.New("UserId is required")
	}

	tx, err := p.prRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = p.userRepo.GetUserByID(ctx, tx, req.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	prEntities, err := p.prRepo.GetPullRequestsByReviewer(ctx, tx, req.UserId)
	if err != nil {
		return nil, err
	}

	pullRequests := make([]*models.PullRequestShort, 0, len(prEntities))
	for _, prEntity := range prEntities {
		pr := &models.PullRequestShort{
			PullRequestId:   prEntity.PullRequestID,
			PullRequestName: prEntity.PullRequestName,
			AuthorId:        prEntity.AuthorID,
			Status:          api.PullRequestShortStatus(prEntity.Status),
		}
		pullRequests = append(pullRequests, pr)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return pullRequests, nil
}
