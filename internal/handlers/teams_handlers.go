package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	api "github.com/oooooorg/PR-Service/internal/gen"
	"github.com/oooooorg/PR-Service/internal/service"
)

func (s *Server) PostTeamAdd(ctx echo.Context) error {
	var body api.PostTeamAddJSONRequestBody

	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	team, err := s.teamService.CreateTeam(ctx.Request().Context(), &body)
	if err != nil {
		if errors.Is(err, service.ErrTeamExists) {
			return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
				Error: struct {
					Code    api.ErrorResponseErrorCode `json:"code"`
					Message string                     `json:"message"`
				}{
					Code:    api.TEAMEXISTS,
					Message: body.TeamName + "already exists",
				},
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, team)
}

func (s *Server) GetTeamGet(ctx echo.Context, params api.GetTeamGetParams) error {
	team, err := s.teamService.GetTeam(ctx.Request().Context(), &params)
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

	return ctx.JSON(http.StatusOK, team)
}
