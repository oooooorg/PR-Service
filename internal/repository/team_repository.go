package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/oooooorg/PR-Service/internal/entity"
)

type TeamRepositoryImpl struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepositoryImpl {
	return &TeamRepositoryImpl{
		db: db,
	}
}

func (tr *TeamRepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return tr.db.BeginTx(ctx, nil)
}

func (tr *TeamRepositoryImpl) CreateTeam(ctx context.Context, tx *sql.Tx, team *entity.Team) error {
	const query = `INSERT INTO teams (team_name) VALUES ($1) RETURNING id`

	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, team.TeamName).Scan(&team.ID)
	} else {
		err = tr.db.QueryRowContext(ctx, query, team.TeamName).Scan(&team.ID)
	}

	return err
}

func (tr *TeamRepositoryImpl) GetTeamByName(ctx context.Context, tx *sql.Tx, teamName string) (*entity.Team, error) {
	const query = `SELECT id, team_name FROM teams WHERE team_name = $1`

	var team entity.Team
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, teamName).Scan(&team.ID, &team.TeamName)
	} else {
		err = tr.db.QueryRowContext(ctx, query, teamName).Scan(&team.ID, &team.TeamName)
	}

	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (tr *TeamRepositoryImpl) TeamExists(ctx context.Context, tx *sql.Tx, teamName string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)`

	var exists bool
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, teamName).Scan(&exists)
	} else {
		err = tr.db.QueryRowContext(ctx, query, teamName).Scan(&exists)
	}

	if err != nil {
		return false, fmt.Errorf("failed to check team existence: %w", err)
	}

	return exists, nil
}
