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

	"github.com/oooooorg/PR-Service/internal/handlers"
	"github.com/oooooorg/PR-Service/internal/handlers/mocks"
	"github.com/oooooorg/PR-Service/internal/models"
)

func newTestServerTeam(teamSerivceMock *mocks.MockTeamService) *handlers.Server {
	return &handlers.Server{
		TeamService: teamSerivceMock,
	}
}

func TestPostTeamAdd_Success(t *testing.T) {
	e := echo.New()

	body := `{
        "team_name": "backend",
        "members": [
            { "user_id": "u1", "username": "Alice", "is_active": true },
            { "user_id": "u2", "username": "Bob", "is_active": true }
        ]
    }`

	request := httptest.NewRequest(http.MethodPost, "/team/add", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)

	teamSerivceMock := new(mocks.MockTeamService)

	teamSerivceMock.
		On(
			"CreateTeam",
			mock.Anything,
			mock.AnythingOfType("*api.Team"),
		).
		Return(
			&models.Team{
				TeamName: "backend",
				Members:  []models.TeamMember{},
			},
			nil,
		)

	serverMock := newTestServerTeam(teamSerivceMock)

	err := serverMock.PostTeamAdd(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, recorder.Code)
	teamSerivceMock.AssertExpectations(t)
}

func TestPostTeamAdd_AlreadyExists(t *testing.T) {
	e := echo.New()

	body := `{
        "team_name": "backend",
        "members": []
    }`

	request := httptest.NewRequest(http.MethodPost, "/team/add", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)

	teamSerivceMock := new(mocks.MockTeamService)

	teamSerivceMock.
		On(
			"CreateTeam",
			mock.Anything,
			mock.AnythingOfType("*api.Team"),
		).
		Return(
			(*models.Team)(nil),
			service.ErrTeamExists,
		)

	serverMock := newTestServerTeam(teamSerivceMock)

	err := serverMock.PostTeamAdd(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	teamSerivceMock.AssertExpectations(t)
}
