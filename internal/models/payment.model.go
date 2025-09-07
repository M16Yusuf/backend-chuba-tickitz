package models

type Payment struct {
	Id     int    `db:"id" json:"payment_id"`
	Method string `db:"method" json:"payment_method"`
}
