package repository

import (
	"context"
	"database/sql"

	"github.com/oooooorg/PR-Service/internal/database"
	"github.com/oooooorg/PR-Service/internal/entity"
)

type PullRequestRepositoryImpl struct {
	db *sql.DB
}

func NewPullRequestRepository(db *sql.DB) *PullRequestRepositoryImpl {
	return &PullRequestRepositoryImpl{
		db: db,
	}
}

func (ps *PullRequestRepositoryImpl) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return ps.db.BeginTx(ctx, nil)
}

func (ps *PullRequestRepositoryImpl) CreatePullRequest(ctx context.Context, tx *sql.Tx, pr *entity.PullRequest) error {
	const query = `
        INSERT INTO pull_requests (
            author_id, pull_request_id, pull_request_name, 
            assigned_reviewers_first, assigned_reviewers_second, status,
            created_at, updated_at_utc
        ) 
        VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) 
        RETURNING id, created_at, updated_at_utc`

	args := []any{
		pr.AuthorID,
		pr.PullRequestID,
		pr.PullRequestName,
		database.StringToNullString(pr.AssignedReviewersFirst),
		database.StringToNullString(pr.AssignedReviewersSecond),
		pr.Status,
	}

	if tx != nil {
		return tx.QueryRowContext(ctx, query, args...).Scan(&pr.ID, &pr.CreatedAt, &pr.UpdatedAt)
	}
	return ps.db.QueryRowContext(ctx, query, args...).Scan(&pr.ID, &pr.CreatedAt, &pr.UpdatedAt)
}

func (ps *PullRequestRepositoryImpl) UpdatePullRequestReviewers(ctx context.Context, tx *sql.Tx, prID string, reviewer1, reviewer2 string) (*entity.PullRequest, error) {
	const query = `
		UPDATE pull_requests 
        SET assigned_reviewers_first = $1, 
            assigned_reviewers_second = $2, 
            updated_at_utc = NOW() 
        WHERE pull_request_id = $3
        RETURNING id, author_id, pull_request_id, pull_request_name, 
                  assigned_reviewers_first, assigned_reviewers_second, 
                  status, created_at, updated_at_utc, merged_at`
	args := []any{
		database.StringToNullString(reviewer1),
		database.StringToNullString(reviewer2),
		prID,
	}

	var pr entity.PullRequest
	var err error
	var rev1, rev2 sql.NullString

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, args...).Scan(
			&pr.ID, &pr.AuthorID, &pr.PullRequestID, &pr.PullRequestName,
			&rev1, &rev2,
			&pr.Status, &pr.CreatedAt, &pr.UpdatedAt, &pr.MergedAt,
		)
	} else {
		err = ps.db.QueryRowContext(ctx, query, args...).Scan(
			&pr.ID, &pr.AuthorID, &pr.PullRequestID, &pr.PullRequestName,
			&rev1, &rev2,
			&pr.Status, &pr.CreatedAt, &pr.UpdatedAt, &pr.MergedAt,
		)
	}

	if err != nil {
		return nil, err
	}

	if rev1.Valid {
		pr.AssignedReviewersFirst = rev1.String
	}
	if rev2.Valid {
		pr.AssignedReviewersSecond = rev2.String
	}

	return &pr, nil
}

func (ps *PullRequestRepositoryImpl) UpdatePullRequestStatus(ctx context.Context, tx *sql.Tx, prID string, status string) (*entity.PullRequest, error) {
	var query string
	if status == "MERGED" {
		query = `
            UPDATE pull_requests 
            SET status = $1, updated_at_utc = NOW(), merged_at = NOW()
            WHERE pull_request_id = $2
            RETURNING id, author_id, pull_request_id, pull_request_name, 
                      assigned_reviewers_first, assigned_reviewers_second, 
                      status, created_at, updated_at_utc, merged_at`
	}

	args := []any{status, prID}

	var pr entity.PullRequest
	var err error
	var rev1, rev2 sql.NullString

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, args...).Scan(
			&pr.ID, &pr.AuthorID, &pr.PullRequestID, &pr.PullRequestName,
			&rev1, &rev2,
			&pr.Status, &pr.CreatedAt, &pr.UpdatedAt, &pr.MergedAt,
		)
	} else {
		err = ps.db.QueryRowContext(ctx, query, args...).Scan(
			&pr.ID, &pr.AuthorID, &pr.PullRequestID, &pr.PullRequestName,
			&rev1, &rev2,
			&pr.Status, &pr.CreatedAt, &pr.UpdatedAt, &pr.MergedAt,
		)
	}

	if err != nil {
		return nil, err
	}

	if rev1.Valid {
		pr.AssignedReviewersFirst = rev1.String
	}
	if rev2.Valid {
		pr.AssignedReviewersSecond = rev2.String
	}

	return &pr, nil
}

func (ps *PullRequestRepositoryImpl) GetPullRequestByID(ctx context.Context, tx *sql.Tx, prID string) (*entity.PullRequest, error) {
	const query = `
        SELECT id, author_id, pull_request_id, pull_request_name, 
               assigned_reviewers_first, assigned_reviewers_second, 
               status, created_at, updated_at_utc, merged_at
        FROM pull_requests
        WHERE pull_request_id = $1
    `

	var pr entity.PullRequest
	var err error
	var rev1, rev2 sql.NullString

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, prID).Scan(
			&pr.ID, &pr.AuthorID, &pr.PullRequestID, &pr.PullRequestName,
			&rev1, &rev2,
			&pr.Status, &pr.CreatedAt, &pr.UpdatedAt, &pr.MergedAt,
		)
	} else {
		err = ps.db.QueryRowContext(ctx, query, prID).Scan(
			&pr.ID, &pr.AuthorID, &pr.PullRequestID, &pr.PullRequestName,
			&rev1, &rev2,
			&pr.Status, &pr.CreatedAt, &pr.UpdatedAt, &pr.MergedAt,
		)
	}

	if err != nil {
		return nil, err
	}

	if rev1.Valid {
		pr.AssignedReviewersFirst = rev1.String
	}
	if rev2.Valid {
		pr.AssignedReviewersSecond = rev2.String
	}

	return &pr, nil
}

func (ps *PullRequestRepositoryImpl) GetPullRequestsByReviewer(ctx context.Context, tx *sql.Tx, reviewerID string) ([]*entity.PullRequest, error) {
	const query = `
        SELECT id, author_id, pull_request_id, pull_request_name, 
               assigned_reviewers_first, assigned_reviewers_second, 
               status, created_at, updated_at_utc, merged_at
        FROM pull_requests
        WHERE assigned_reviewers_first = $1 OR assigned_reviewers_second = $1
        ORDER BY created_at DESC
    `

	var rows *sql.Rows
	var err error

	if tx != nil {
		rows, err = tx.QueryContext(ctx, query, reviewerID)
	} else {
		rows, err = ps.db.QueryContext(ctx, query, reviewerID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pullRequests []*entity.PullRequest
	for rows.Next() {
		var pr entity.PullRequest
		var rev1, rev2 sql.NullString

		err := rows.Scan(
			&pr.ID, &pr.AuthorID, &pr.PullRequestID, &pr.PullRequestName,
			&rev1, &rev2,
			&pr.Status, &pr.CreatedAt, &pr.UpdatedAt, &pr.MergedAt,
		)
		if err != nil {
			return nil, err
		}

		if rev1.Valid {
			pr.AssignedReviewersFirst = rev1.String
		}
		if rev2.Valid {
			pr.AssignedReviewersSecond = rev2.String
		}

		pullRequests = append(pullRequests, &pr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pullRequests, nil
}
