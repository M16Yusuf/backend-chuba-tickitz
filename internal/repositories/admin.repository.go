package repositories

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
)

type AdminRepository struct {
	db *pgxpool.Pool
}

func NewAdminRepository(db *pgxpool.Pool) *AdminRepository {
	return &AdminRepository{db: db}
}

// Get all movie for admih dashboard page
// Get data movie with page (limit 20, offset)
// Query table (movies, movies_genres)
func (a *AdminRepository) GetAllMovies(reqCntxt context.Context, offset, limit int) ([]models.MovieList, error) {
	sql := `SELECT m.id, m.title,  m.poster_path, m.release_date, m.duration, 
		json_agg(json_build_object('genre_id', g.id, 'genre_name', g.name)) AS genres
		FROM movies m
		LEFT JOIN genres_movies gm ON m.id = gm.movie_id
		LEFT JOIN genres g ON gm.genre_id = g.id
		GROUP BY m.id, m.title, m.poster_path, m.release_date, m.duration
		ORDER BY m.release_date DESC
		LIMIT $2 OFFSET $1`

	values := []any{offset, limit}
	rows, err := a.db.Query(reqCntxt, sql, values...)
	if err != nil {
		log.Println("internal server error : ", err.Error())
		return []models.MovieList{}, err
	}
	defer rows.Close()

	// processing data / read rows
	var movies []models.MovieList
	for rows.Next() {
		var movie models.MovieList
		var genreRaw []byte
		if err := rows.Scan(&movie.Id, &movie.Title, &movie.Poster, &movie.Release_date, &movie.Duration, &genreRaw); err != nil {
			log.Println("Scan Error, ", err.Error())
			return []models.MovieList{}, err
		}

		// decode raw JSOn into []Genres
		if err := json.Unmarshal(genreRaw, &movie.Genres); err != nil {
			log.Println("Unmarshal Error:", err)
			return nil, err
		}

		movies = append(movies, movie)
	}

	return movies, nil
}
