package handlers_test

import (
	"github.com/oooooorg/PR-Service/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	api "github.com/oooooorg/PR-Service/internal/gen"
	"github.com/oooooorg/PR-Service/internal/handlers"
	"github.com/oooooorg/PR-Service/internal/handlers/mocks"
	"github.com/oooooorg/PR-Service/internal/models"
)

func newTestServerPullRequest(pullRequestServiceMock *mocks.MockPullRequestService) *handlers.Server {
	return &handlers.Server{
		PullRequestService: pullRequestServiceMock,
	}
}

func TestPostPullRequestCreate_Success(t *testing.T) {
	e := echo.New()

	body := `{
        "pull_request_id": "pr-1001",
        "pull_request_name": "Add search",
        "author_id": "u1"
    }`

	request := httptest.NewRequest(http.MethodPost, "/pullRequest/create", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)

	pullRequestServiceMock := new(mocks.MockPullRequestService)

	pullRequestServiceMock.
		On(
			"CreatePullRequest",
			mock.Anything,
			mock.AnythingOfType("*api.PostPullRequestCreateJSONRequestBody"),
		).
		Return(
			&models.PullRequest{
				PullRequestId:     "pr-1001",
				PullRequestName:   "Add search",
				AuthorId:          "u1",
				Status:            api.PullRequestStatusOPEN,
				AssignedReviewers: []string{"u2", "u3"},
			},
			nil,
		)

	serverMock := newTestServerPullRequest(pullRequestServiceMock)

	err := serverMock.PostPullRequestCreate(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, recorder.Code)
	pullRequestServiceMock.AssertExpectations(t)
}

func TestPostPullRequestCreate_Conflict(t *testing.T) {
	e := echo.New()

	body := `{
        "pull_request_id": "pr-1001",
        "pull_request_name": "Add search",
        "author_id": "u1"
    }`

	request := httptest.NewRequest(http.MethodPost, "/pullRequest/create", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)

	pullRequestServiceMock := new(mocks.MockPullRequestService)

	pullRequestServiceMock.
		On(
			"CreatePullRequest",
			mock.Anything,
			mock.AnythingOfType("*api.PostPullRequestCreateJSONRequestBody"),
		).
		Return(
			(*models.PullRequest)(nil),
			service.ErrPullRequestExists,
		)

	serverMock := newTestServerPullRequest(pullRequestServiceMock)

	err := serverMock.PostPullRequestCreate(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, recorder.Code)
	pullRequestServiceMock.AssertExpectations(t)
}

func TestPostPullRequestMerge_Success(t *testing.T) {
	e := echo.New()

	body := `{"pull_request_id": "pr-1001"}`

	request := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)

	pullRequestServiceMock := new(mocks.MockPullRequestService)

	pullRequestServiceMock.
		On(
			"MergePullRequest",
			mock.Anything,
			mock.AnythingOfType("*api.PostPullRequestMergeJSONRequestBody"),
		).
		Return(
			&models.PullRequest{
				PullRequestId:     "pr-1001",
				PullRequestName:   "Add search",
				AuthorId:          "u1",
				Status:            api.PullRequestStatusMERGED,
				AssignedReviewers: []string{"u2", "u3"},
			},
			nil,
		)

	serverMock := newTestServerPullRequest(pullRequestServiceMock)

	err := serverMock.PostPullRequestMerge(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	pullRequestServiceMock.AssertExpectations(t)
}

func TestPostPullRequestReassign_Success(t *testing.T) {
	e := echo.New()

	body := `{
        "pull_request_id": "pr-1001",
        "old_user_id": "u2"
    }`

	request := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)

	pullRequestServiceMock := new(mocks.MockPullRequestService)

	pullRequestServiceMock.
		On(
			"ReassignReviewer",
			mock.Anything,
			mock.Anything,
		).
		Return(
			&models.PullRequest{
				PullRequestId:     "pr-1001",
				PullRequestName:   "Add search",
				AuthorId:          "u1",
				Status:            api.PullRequestStatusOPEN,
				AssignedReviewers: []string{"u3", "u5"},
			},
			"u5",
			nil,
		)

	serverMock := newTestServerPullRequest(pullRequestServiceMock)

	err := serverMock.PostPullRequestReassign(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	pullRequestServiceMock.AssertExpectations(t)
}
