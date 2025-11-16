package repository

import (
	"context"
	"database/sql"

	"github.com/oooooorg/PR-Service/internal/entity"
)

type UserRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (ur *UserRepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return ur.db.BeginTx(ctx, nil)
}

func (ur *UserRepositoryImpl) CreateUser(ctx context.Context, tx *sql.Tx, user *entity.User) error {
	const query = `INSERT INTO users (user_id, username, is_active, team_name) VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	args := []any{user.UserID, user.Username, user.IsActive, user.TeamName}

	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
	}
	return ur.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt)
}

func (ur *UserRepositoryImpl) SetUserActive(ctx context.Context, tx *sql.Tx, userID string, isActive bool) (*entity.User, error) {
	const query = `UPDATE users SET is_active = $1, updated_at = NOW() WHERE user_id = $2 RETURNING id, user_id, username, team_name, is_active, created_at, updated_at`

	args := []any{isActive, userID}

	var user entity.User
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, args...).Scan(
			&user.ID, &user.UserID, &user.Username, &user.TeamName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		err = ur.db.QueryRowContext(ctx, query, args...).Scan(
			&user.ID, &user.UserID, &user.Username, &user.TeamName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepositoryImpl) GetUsersByTeam(ctx context.Context, tx *sql.Tx, teamName string) ([]*entity.User, error) {
	const query = `
        SELECT id, user_id, username, team_name, is_active, created_at, updated_at
        FROM users
        WHERE team_name = $1
    `

	args := []any{teamName}

	var rows *sql.Rows
	var err error

	if tx != nil {
		rows, err = tx.QueryContext(ctx, query, args...)
	} else {
		rows, err = ur.db.QueryContext(ctx, query, args...)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.UserID, &user.Username, &user.TeamName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *UserRepositoryImpl) GetUserByID(ctx context.Context, tx *sql.Tx, userID string) (*entity.User, error) {
	const query = `
        SELECT id, user_id, username, team_name, is_active, created_at, updated_at
        FROM users
        WHERE user_id = $1
    `

	var user entity.User
	var err error

	args := []any{userID}

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, args...).Scan(
			&user.ID, &user.UserID, &user.Username, &user.TeamName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
	} else {
		err = ur.db.QueryRowContext(ctx, query, args...).Scan(
			&user.ID, &user.UserID, &user.Username, &user.TeamName, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
