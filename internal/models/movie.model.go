package models

import "time"

type Movie struct {
	Id           int        `db:"id" json:"movie_id"`
	Poster       *string    `db:"poster_path" json:"poster_path,omitempty"`
	Backdrop     *string    `db:"backdrop_path" json:"backdrop_path,omitempty"`
	Title        *string    `db:"title" json:"title,omitempty"`
	Overview     *string    `db:"overview" json:"overview,omitempty"`
	Release_date *time.Time `db:"release_date" json:"release_date,omitempty"`
	Duration     *int       `db:"duration" json:"duration,omitempty"`
	Director_id  *int       `db:"director_id" json:"director_id,omitempty"`
}

type MovieList struct {
	Movie
	Genres []Genre
}
