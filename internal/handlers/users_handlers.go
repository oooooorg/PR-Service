package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	api "github.com/oooooorg/PR-Service/internal/gen"
)

func (s *Server) PostUsersSetIsActive(ctx echo.Context) error {
	var body api.PostUsersSetIsActiveJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := s.userService.SetUserActive(ctx.Request().Context(), &body)
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

	return ctx.JSON(http.StatusOK, user)
}
