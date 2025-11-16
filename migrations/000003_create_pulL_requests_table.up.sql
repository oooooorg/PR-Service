CREATE TABLE IF NOT EXISTS pull_requests (
    id SERIAL PRIMARY KEY,
    author_id VARCHAR(100) NOT NULL REFERENCES users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    merged_at TIMESTAMP,
    pull_request_id VARCHAR(100) UNIQUE NOT NULL,
    pull_request_name TEXT NOT NULL,
    assigned_reviewers_first VARCHAR(100) REFERENCES users(user_id),
    assigned_reviewers_second VARCHAR(100) REFERENCES users(user_id),
    status VARCHAR(50) NOT NULL,
    updated_at_utc TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_pull_requests_pull_request_id ON pull_requests(pull_request_id);
CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_pull_requests_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_pull_requests_created_at ON pull_requests(created_at);
CREATE INDEX IF NOT EXISTS idx_pull_requests_merged_at ON pull_requests(merged_at);
CREATE INDEX IF NOT EXISTS idx_pull_requests_assigned_reviewers_first ON pull_requests(assigned_reviewers_first);
CREATE INDEX IF NOT EXISTS idx_pull_requests_assigned_reviewers_second ON pull_requests(assigned_reviewers_second);
