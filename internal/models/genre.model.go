package models

type Genre struct {
	Id   int    `db:"id" json:"genre_id"`
	Name string `db:"name" json:"genre_name"`
}
