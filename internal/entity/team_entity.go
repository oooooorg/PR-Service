package entity

type Team struct {
	ID       int    `db:"id"`
	TeamName string `db:"team_name"`
}
