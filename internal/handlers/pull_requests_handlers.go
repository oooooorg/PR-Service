package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	api "github.com/oooooorg/PR-Service/internal/gen"
	"github.com/oooooorg/PR-Service/internal/service"
)

func (s *Server) PostPullRequestCreate(ctx echo.Context) error {
	var body api.PostPullRequestCreateJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pr, err := s.pullRequestService.CreatePullRequest(ctx.Request().Context(), &body)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrTeamNotFound) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Error: struct {
					Code    api.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{
					Code:    api.NOTFOUND,
					Message: "resource not found",
				},
			})
		}

		if errors.Is(err, service.ErrPullRequestExists) {
			return ctx.JSON(http.StatusConflict, api.ErrorResponse{
				Error: struct {
					Code    api.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{
					Code:    api.PREXISTS,
					Message: body.PullRequestId + "already exists",
				},
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, pr)
}

func (s *Server) PostPullRequestMerge(ctx echo.Context) error {
	var body api.PostPullRequestMergeJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pr, err := s.pullRequestService.MergePullRequest(ctx.Request().Context(), &body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Error: struct {
					Code    api.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{
					Code:    api.NOTFOUND,
					Message: "resource not found",
				},
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, pr)
}

func (s *Server) PostPullRequestReassign(ctx echo.Context) error {
	var body api.PostPullRequestReassignJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pr, replacedBy, err := s.pullRequestService.ReassignReviewer(ctx.Request().Context(), &body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, service.ErrUserNotFound) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Error: struct {
					Code    api.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{
					Code:    api.NOTFOUND,
					Message: "resource not found",
				},
			})
		}

		if errors.Is(err, service.ErrPullRequestMerged) {
			return ctx.JSON(http.StatusConflict, api.ErrorResponse{
				Error: struct {
					Code    api.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{
					Code:    api.PRMERGED,
					Message: "cannot reassign on merged PR",
				},
			})
		}

		if errors.Is(err, service.ErrPullRequestNotAsigned) {
			return ctx.JSON(http.StatusConflict, api.ErrorResponse{
				Error: struct {
					Code    api.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{
					Code:    api.NOTASSIGNED,
					Message: "reviewer is not assigned to this PR",
				},
			})
		}

		if errors.Is(err, service.ErrPullRequestNoCandidate) {
			return ctx.JSON(http.StatusConflict, api.ErrorResponse{
				Error: struct {
					Code    api.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{
					Code:    api.NOCANDIDATE,
					Message: "no active replacement candidate in team",
				},
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"pr":          pr,
		"replaced_by": replacedBy,
	})
}

func (s *Server) GetUsersGetReview(ctx echo.Context, params api.GetUsersGetReviewParams) error {
	pullRequests, err := s.pullRequestService.GetUserReviewRequests(ctx.Request().Context(), &params)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return ctx.JSON(http.StatusNotFound, api.ErrorResponse{
				Error: struct {
					Code    api.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{
					Code:    api.NOTFOUND,
					Message: "resource not found",
				},
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"user_id":       params.UserId,
		"pull_requests": pullRequests,
	})
}
