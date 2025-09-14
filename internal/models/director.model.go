package models

type Director struct {
	Id   *int    `db:"id" json:"id"`
	Name *string `db:"name" json:"name"`
}
