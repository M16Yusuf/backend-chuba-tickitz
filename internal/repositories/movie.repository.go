package repositories

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/utils"
	"github.com/redis/go-redis/v9"
)

type MovieRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewMovieRepository(db *pgxpool.Pool, rdb *redis.Client) *MovieRepository {
	return &MovieRepository{db: db, rdb: rdb}
}

// Function Query untuk Upcoming movies
// upcoming, movie yang release datenya masih in the future
func (m *MovieRepository) GetUpcoming(reqCntxt context.Context, offset, limit int) ([]models.MovieList, error) {

	// Get cached upcoming movies, before accesing database
	// Get and renew cache only for page 1 only (offset 0)
	rdbKey := "chuba_tickitz:movies-upcoming"
	if offset == 0 {
		cachedUpcomMovies, err := utils.RedisGetData[[]models.MovieList](reqCntxt, *m.rdb, rdbKey)
		if err != nil {
			log.Println("Redis error :", err)
		} else if cachedUpcomMovies != nil && len(*cachedUpcomMovies) > 0 {
			return *cachedUpcomMovies, nil
		}
	}

	// if there is no key/ no cached data, get upcoming movies from database
	sql := `SELECT m.id, m.poster_path, m.title, m.release_date, 
		json_agg(json_build_object('genre_id', g.id, 'genre_name', g.name)) AS genres
		FROM movies m
		JOIN movie_genres mg ON m.id = mg.movie_id
		JOIN genres g ON mg.genre_id = g.id
		WHERE m.release_date > CURRENT_DATE AND m.deleted_at IS NULL
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

	// make cache upcoming movies after query data from database
	// Get and renew cache only for page 1 only (offset 0)
	if offset == 0 {
		if err := utils.RedisRenewData(reqCntxt, *m.rdb, rdbKey, Movies, 5*time.Minute); err != nil {
			log.Println("Failed to renew Redis cache:", err.Error())
		}
	}
	// return data movies ([]model,movielist), and return error nil of not error
	return Movies, nil
}

// Function query for Popular movies
// Popular, movie sorted by rating from transaction
func (m *MovieRepository) GetPopular(reqCntxt context.Context, offset, limit int) ([]models.MovieList, error) {

	// Get cached popular movies, before accesing database
	// Get and renew cache only for page 1 only (offset 0)
	rdbKey := "chuba_tickitz:movies-popular"
	if offset == 0 {
		cachedPopMovies, err := utils.RedisGetData[[]models.MovieList](reqCntxt, *m.rdb, rdbKey)
		if err != nil {
			log.Println("Redis error :", err)
		} else if cachedPopMovies != nil && len(*cachedPopMovies) > 0 {
			return *cachedPopMovies, nil
		}
	}

	// if there is no key/ no cached data, get popular movies from database
	sql := `SELECT m.id, m.poster_path, m.title, AVG(t.rating) AS avg_rating, COUNT(t.id) AS rating_count,
    json_agg(DISTINCT jsonb_build_object('genre_id', g.id, 'genre_name', g.name)) AS genres
		FROM movies m
		JOIN schedules s ON s.movie_id = m.id
		JOIN transactions t ON t.schedule_id = s.id
		JOIN movie_genres mg ON m.id = mg.movie_id
		JOIN genres g ON mg.genre_id = g.id
		WHERE  t.rating IS NOT NULL AND m.deleted_at IS NULL
		GROUP BY m.id, m.poster_path, m.title
		ORDER BY avg_rating DESC, rating_count DESC
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
		if err := rows.Scan(&Movie.Id, &Movie.Poster, &Movie.Title, &Movie.AvgRating, &Movie.RatingCount, &genreRaw); err != nil {
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

	// make cache popular movies after query data from database
	// Get and renew cache only for page 1 only (offset 0)
	if offset == 0 {
		if err := utils.RedisRenewData(reqCntxt, *m.rdb, rdbKey, Movies, 5*time.Minute); err != nil {
			log.Println("Failed to renew Redis cache:", err.Error())
		}
	}
	// return data movies ([]model,movielist), and return error nil of not error
	return Movies, nil
}

// Function query for filter movies
// Filter by title/search and genres
func (m *MovieRepository) GetFiltered(reqCntxt context.Context, offset, limit int, search string, genrSearch []string) ([]models.MovieList, error) {
	// process the query based on what user inputs
	values := []any{}
	sql := `SELECT m.id, m.poster_path, m.title, 
	json_agg(DISTINCT jsonb_build_object('genre_id', g.id, 'genre_name', g.name)) AS genres	
	FROM movies m
	JOIN movie_genres mg ON m.id = mg.movie_id
	JOIN genres g ON mg.genre_id = g.id `
	if search != "" {
		idx := strconv.Itoa(len(values) + 1)
		sql += "WHERE m.title ILIKE '%' || $" + idx + " || '%'	"
		values = append(values, search)
	}
	sql += `GROUP BY m.id, m.poster_path, m.title `
	if len(genrSearch) > 0 {
		idx := strconv.Itoa(len(values) + 1)
		sql += "HAVING ARRAY_AGG(DISTINCT g.name)::text[] @> $" + idx + " "
		values = append(values, pq.Array(genrSearch))
	}
	idx1 := strconv.Itoa(len(values) + 1)
	idx2 := strconv.Itoa(len(values) + 2)
	sql += "LIMIT $" + idx1 + " OFFSET $" + idx2 + " "
	values = append(values, limit)
	values = append(values, offset)

	rows, err := m.db.Query(reqCntxt, sql, values...)
	// log.Println(rows)
	if err != nil {
		log.Println("internal server error : ", err.Error())
		return []models.MovieList{}, err
	}
	defer rows.Close()
	// 	processing data pgx.rows / read rows, append into slice
	var Movies []models.MovieList
	for rows.Next() {
		var Movie models.MovieList
		var genreRaw []byte
		if err := rows.Scan(&Movie.Id, &Movie.Poster, &Movie.Title, &genreRaw); err != nil {
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
	// log.Println("test", Movies)
	return Movies, nil
}

// Function query get movie details
func (m *MovieRepository) GetMovieDetails(reqCntxt context.Context, movieID string) (models.MovieDetails, error) {
	sql := `SELECT m.id, m.title, m.overview, m.poster_path, m.backdrop_path, m.release_date, m.duration, d.name AS director_name, d.id AS director_id, 
	json_agg(DISTINCT jsonb_build_object('genre_id', g.id, 'genre_name', g.name)) AS genres, 
	json_agg(DISTINCT jsonb_build_object('actor_id', a.id, 'actor_name', a.name)) AS actors
  FROM movies m
  JOIN directors d ON m.director_id = d.id
	LEFT JOIN movie_genres mg ON m.id = mg.movie_id
	LEFT JOIN genres g ON mg.genre_id = g.id
  LEFT JOIN movie_actors ma ON m.id = ma.movie_id
  LEFT JOIN actors a ON ma.actor_id = a.id
  WHERE m.id = $1
  GROUP BY m.id, m.title, m.overview, m.poster_path, m.backdrop_path, m.release_date, m.duration, d.name, d.id;`

	var MovieDetails models.MovieDetails
	var genreRaw []byte
	var actorRaw []byte
	err := m.db.QueryRow(reqCntxt, sql, movieID).Scan(&MovieDetails.Id, &MovieDetails.Title, &MovieDetails.Overview, &MovieDetails.Poster, &MovieDetails.Backdrop, &MovieDetails.Release_date, &MovieDetails.Duration, &MovieDetails.Director.Name, &MovieDetails.Director.Id, &genreRaw, &actorRaw)
	if err != nil {
		log.Println("Scan Error, ", err.Error())
		return models.MovieDetails{}, err
	}

	// Decode raw JSON into []Genre
	if len(genreRaw) > 0 {
		if err := json.Unmarshal(genreRaw, &MovieDetails.Genres); err != nil {
			log.Println("Unmarshal Error:", err)
			return models.MovieDetails{}, err
		}
	}

	// Decode raw JSON into []Actor
	if len(actorRaw) > 0 {
		if err := json.Unmarshal(actorRaw, &MovieDetails.Actors); err != nil {
			log.Println("Unmarshal Error:", err)
			return models.MovieDetails{}, err
		}
	}

	return MovieDetails, nil
}

// function query get all genres
func (m *MovieRepository) GetGenres(rqCntxt context.Context) ([]models.Genre, error) {

	// Get cached genres from redis, before accesing database
	rdbKey := "chuba_tickitz:all-genres"
	cachedGenres, err := utils.RedisGetData[[]models.Genre](rqCntxt, *m.rdb, rdbKey)
	if err != nil {
		log.Println("Redis error :", err)
	} else if cachedGenres != nil && len(*cachedGenres) > 0 {
		return *cachedGenres, nil
	}

	// if there is no key/ no cached data, get all genres from database
	sql := `SELECT id, name FROM genres`
	rows, err := m.db.Query(rqCntxt, sql)
	if err != nil {
		log.Println("internal server error : ", err.Error())
		return []models.Genre{}, err
	}
	defer rows.Close()

	// process the rows
	var genres []models.Genre
	for rows.Next() {
		var genre models.Genre
		if err := rows.Scan(&genre.Id, &genre.Name); err != nil {
			log.Println("Scan Error, ", err.Error())
			return []models.Genre{}, err
		}
		genres = append(genres, genre)
	}

	// make cache genre after query data from database
	if err := utils.RedisRenewData(rqCntxt, *m.rdb, rdbKey, genres, 5*time.Minute); err != nil {
		log.Println("Failed to renew Redis cache:", err.Error())
	}

	return genres, nil
}
