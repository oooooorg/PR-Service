CREATE TABLE IF NOT EXISTS teams
(
    id SERIAL PRIMARY KEY,
    team_name VARCHAR(100) UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_teams_team_name_idx ON teams(team_name);
