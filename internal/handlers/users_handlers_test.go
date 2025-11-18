package handlers_test

import (
	"database/sql"
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

func newTestServerUser(userServiceMock *mocks.MockUserService) *handlers.Server {
	return &handlers.Server{
		UserService: userServiceMock,
	}
}

func TestPostUsersSetIsActive_Success(t *testing.T) {
	e := echo.New()

	body := `{
        "user_id": "u2",
        "is_active": false
    }`

	request := httptest.NewRequest(http.MethodPost, "/users/setIsActive", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)

	userServiceMock := new(mocks.MockUserService)

	userServiceMock.
		On(
			"SetUserActive",
			mock.Anything,
			mock.AnythingOfType("*api.PostUsersSetIsActiveJSONRequestBody"),
		).
		Return(
			&models.User{
				UserId:   "u2",
				Username: "Bob",
				TeamName: "backend",
				IsActive: false,
			},
			nil,
		)

	serverMock := newTestServerUser(userServiceMock)

	err := serverMock.PostUsersSetIsActive(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, recorder.Code)
	userServiceMock.AssertExpectations(t)
}

func TestPostUsersSetIsActive_NotFound(t *testing.T) {
	e := echo.New()

	body := `{
        "user_id": "CTitmo",
        "is_active": true
    }`

	request := httptest.NewRequest(http.MethodPost, "/users/setIsActive", strings.NewReader(body))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	recorder := httptest.NewRecorder()
	ctx := e.NewContext(request, recorder)

	userServiceMock := new(mocks.MockUserService)

	userServiceMock.
		On(
			"SetUserActive",
			mock.Anything,
			mock.AnythingOfType("*api.PostUsersSetIsActiveJSONRequestBody"),
		).
		Return(
			(*models.User)(nil),
			sql.ErrNoRows,
		)

	serverMock := newTestServerUser(userServiceMock)

	err := serverMock.PostUsersSetIsActive(ctx)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	userServiceMock.AssertExpectations(t)
}
