package models

type City struct {
	Id   *int    `db:"id" json:"city_id"`
	Name *string `db:"name" json:"city_name"`
}
