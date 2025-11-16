DROP INDEX IF EXISTS idx_pull_requests_pull_request_id;
DROP INDEX IF EXISTS idx_pull_requests_author_id;
DROP INDEX IF EXISTS idx_pull_requests_status;
DROP INDEX IF EXISTS idx_pull_requests_created_at;
DROP INDEX IF EXISTS idx_pull_requests_merged_at;
DROP INDEX IF EXISTS idx_pull_requests_assigned_reviewers_first;
DROP INDEX IF EXISTS idx_pull_requests_assigned_reviewers_second;

DROP TABLE IF EXISTS pull_requests;
