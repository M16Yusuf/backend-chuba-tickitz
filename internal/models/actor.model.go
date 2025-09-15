package models

type Actor struct {
	Id   int    `db:"id" json:"actor_id"`
	Name string `db:"name" json:"actor_name"`
}
