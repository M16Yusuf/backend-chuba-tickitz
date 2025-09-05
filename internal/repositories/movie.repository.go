package repositories

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{db: db}
}

// Function Query untuk Upcoming movies
// upcoming, movie yang release datenya masih in the future
func (m *MovieRepository) GetUpcoming(reqCntxt context.Context, offset, limit int) ([]models.MovieList, error) {
	sql := `SELECT m.id, m.poster_path, m.title, m.release_date, 
		json_agg(json_build_object('genre_id', g.id, 'genre_name', g.name)) AS genres
		FROM movies m
		JOIN genres_movies gm ON m.id = gm.movie_id
		JOIN genres g ON gm.genre_id = g.id
		WHERE m.release_date > CURRENT_DATE
		GROUP BY m.id, m.poster_path, m.title, m.release_date
		LIMIT $2 OFFSET $1`

	values := []any{offset, limit}
	rows, err := m.db.Query(reqCntxt, sql, values...)
	if err != nil {
		log.Println("internal server error : ", err.Error())
		return []models.MovieList{}, err
	}
	defer rows.Close()

	// 	processing data / read rows
	var Movies []models.MovieList
	for rows.Next() {
		var Movie models.MovieList
		var genreRaw []byte
		if err := rows.Scan(&Movie.Id, &Movie.Poster, &Movie.Title, &Movie.Release_date, &genreRaw); err != nil {
			log.Println("Scan Error, ", err.Error())
			return []models.MovieList{}, err
		}

		// Decode raw JSON into []Genre
		if err := json.Unmarshal(genreRaw, &Movie.Genres); err != nil {
			log.Println("Unmarshal Error:", err)
			return nil, err
		}

		Movies = append(Movies, Movie)
	}
	log.Println(Movies)
	return Movies, nil
}
