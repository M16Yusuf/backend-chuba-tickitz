package models

import "time"

type Transaction struct {
	Id         int       `db:"id" json:"transaction_id"`
	Code       string    `db:"code_ticket" json:"code_ticket"`
	MovieID    int       `db:"movies_id" json:"movie_id"`
	CityID     int       `db:"city_id" json:"city_id"`
	CinemaID   int       `db:"cinema_id" json:"cinema_id"`
	ScheduleID int       `db:"schedule_id" json:"schedule_id"`
	PaymentID  int       `db:"payment_id" json:"payment_id"`
	PaidAt     time.Time `db:"paid_at" json:"paid_at"`
	Total      float64   `db:"total_price" json:"total"`
	CreatedAt  float64   `db:"created_at" json:"created_at"`
	Rating     float64   `db:"rating" json:"rating"`
}

type History struct {
	Transaction
	Movie    Movie
	Schedule Schedule
	Payment  Payment
	Cinema   Cinema
	City     City
	Seat     []Seat
}
