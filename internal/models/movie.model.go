package models

import "time"

type Movie struct {
	Id           int        `db:"id" json:"movie_id"`
	Poster       *string    `db:"poster_path" json:"poster_path"`
	Backdrop     *string    `db:"backdrop_path" json:"backdrop_path"`
	Title        *string    `db:"title" json:"title"`
	Overview     *string    `db:"overview" json:"overview"`
	Release_date *time.Time `db:"release_date" json:"release_date"`
	Duration     *int       `db:"duration" json:"duration"`
	Director_id  *int       `db:"director_id" json:"director_id"`
}

type MovieList struct {
	Movie
	Genres []Genre `json:"genres"`
}

type MovieDetails struct {
	Movie
	Genres   []Genre `json:"genres"`
	Actors   []Actor `json:"actors"`
	Director Director
}
