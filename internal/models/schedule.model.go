package models

import "time"

type Schedule struct {
	Id        int       `db:"id" json:"schedule_id"`
	MovieId   int       `db:"movie_id" json:"movie_id"`
	MovieName string    `db:"title" json:"title"`
	Schedule  time.Time `db:"schedule" json:"schedule"`
}
