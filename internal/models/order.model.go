package models

import "time"

type Order struct {
	CodeTicket *string    `db:"code_ticket" json:"code_ticket"`
	MovieId    int        `db:"movie_id" json:"movie_id"`
	UserId     int        `db:"user_id" json:"user_id"`
	ScheduleId int        `db:"schedule_id" json:"schedule_id"`
	PaymentId  *int       `db:"payment_id" json:"payment_id"`
	PaidAt     *time.Time `db:"paid_at" json:"paid_at"`
	TotalPrice int        `db:"total_price" json:"total_price"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	Rating     *float64   `db:"rating" json:"rating"`
}

type CreateOrder struct {
	Order
	Seats []Seat `db:"seats" json:"seats"`
}
