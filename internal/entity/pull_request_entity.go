package entity

import "time"

type PullRequest struct {
	ID                      int               `db:"id"`
	AuthorID                string            `db:"author_id"`
	PullRequestID           string            `db:"pull_request_id"`
	PullRequestName         string            `db:"pull_request_name"`
	AssignedReviewersFirst  string            `db:"assigned_reviewers_first"`
	AssignedReviewersSecond string            `db:"assigned_reviewers_second"`
	Status                  PullRequestStatus `db:"status"`
	CreatedAt               time.Time         `db:"created_at"`
	MergedAt                *time.Time        `db:"merged_at"`
	UpdatedAt               time.Time         `db:"updated_at"`
}

type PullRequestStatus string

const (
	StatusOpen   PullRequestStatus = "OPEN"
	StatusMerged PullRequestStatus = "MERGED"
)
