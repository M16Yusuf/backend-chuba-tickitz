package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/utils"
	"github.com/redis/go-redis/v9"
)

type AdminRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAdminRepository(db *pgxpool.Pool, rdb *redis.Client) *AdminRepository {
	return &AdminRepository{db: db, rdb: rdb}
}

// Get all movie for admih dashboard page
// Get data movie with page (limit 20, offset)
// Query table (movies, movies_genres)
func (a *AdminRepository) GetAllMovies(reqCntxt context.Context, offset, limit int) ([]models.MovieList, error) {

	// Get cached all movies, before accesing database
	// Get and renew cache only for page 1 only (offset 0)
	rdbKey := "chuba_tickitz:admin-allmovies"
	if offset == 0 {
		cachedAllMovies, err := utils.RedisGetData[[]models.MovieList](reqCntxt, *a.rdb, rdbKey)
		if err != nil {
			log.Println("Redis error :", err)
		} else if cachedAllMovies != nil && len(*cachedAllMovies) > 0 {
			return *cachedAllMovies, nil
		}
	}

	sql := `SELECT m.id, m.title,  m.poster_path, m.release_date, m.duration, 
		json_agg(json_build_object('genre_id', g.id, 'genre_name', g.name)) AS genres
		FROM movies m
		LEFT JOIN movie_genres mg ON m.id = mg.movie_id
		LEFT JOIN genres g ON mg.genre_id = g.id
		WHERE m.deleted_at IS NULL
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

	// make cache all movies after query data from database
	// Get and renew cache only for page 1 only (offset 0)
	if offset == 0 {
		if err := utils.RedisRenewData(reqCntxt, *a.rdb, rdbKey, movies, 10*time.Minute); err != nil {
			log.Println("Failed to renew Redis cache:", err.Error())
		}
	}
	// return data movies ([]model.movielist) , and errror nil if not error
	return movies, nil
}

// Add data new movie, require admin role
// Add data movie, actor, movie_genres, Director
// Query effected table : movie, actor, movie_genres, director
func (a *AdminRepository) AddMovie(reqCntxt context.Context, body models.MovieDetails) error {

	return nil
}

// Delete a movie, require admin role
// soft delete: just set time on column deleted_at
// Query effected table : mmovies only
func (a *AdminRepository) DeleteMovie(reqCntxt context.Context, movieID string) error {
	sql := `UPDATE movies SET deleted_at=CURRENT_TIMESTAMP WHERE id=$1`
	// change query input become int instead of string
	newMovieID, err := strconv.Atoi(movieID)
	if err != nil {
		log.Println("Failed convert string to int\nCause : ", err)
		return err
	}

	cmd, err := a.db.Exec(reqCntxt, sql, newMovieID)
	if err != nil {
		log.Println("Failed execute query\nCause:", err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		log.Println("No movies found with given id: ")
		errormsg := fmt.Sprintf("no movies found with given id: %s", movieID)
		return errors.New(errormsg)
	}

	// if no error return nil
	log.Println("Movie deleted successfully:", movieID)
	return nil
}
