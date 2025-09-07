package models

type Cinema struct {
	Id   int    `db:"id" json:"cinema_id"`
	Name string `db:"name" json:"cinema_name"`
}
