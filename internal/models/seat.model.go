package models

import "time"

type Seat struct {
	Id   int    `db:"id" json:"seat_id"`
	Code string `db:"code" json:"seat_code"`
}

type BookedSeatBySchedule struct {
	ScheduleId    int       `db:"schedule_id" json:"schedule_id"`
	Title         string    `db:"title" json:"title"`
	Schedule      time.Time `db:"schedule" json:"schedule"`
	Code          string    `db:"seat_code" json:"seat_code"`
	TransactionId int       `db:"transaction_id" json:"transaction_id"`
	PaidAt        time.Time `db:"paid_at" json:"paid_at"`
}
