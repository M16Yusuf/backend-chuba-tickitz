package models

import "time"

type Transaction struct {
	Id        int        `db:"id" json:"transaction_id"`
	Code      *string    `db:"code_ticket" json:"code_ticket"`
	PaymentID *int       `db:"payment_id" json:"payment_id"`
	PaidAt    *time.Time `db:"paid_at" json:"paid_at"`
	Total     float64    `db:"total_price" json:"total"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	Rating    *float64   `db:"rating" json:"rating"`
}

type History struct {
	Transaction
	MovieTitle     string    `json:"movie_title"`
	ScheduleTime   time.Time `json:"schedule_time"`
	Payment_method string    `json:"payment_method"`
	CinemaName     string    `json:"cinema_name"`
	Cityname       string    `json:"city_name"`
	Seat           []Seat    `json:"seats"`
}

// id, code_ticket, paid_at, total_price, created_at, rating
