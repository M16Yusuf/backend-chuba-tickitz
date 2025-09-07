package repositories

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
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
	return Movies, nil
}

// Function query for Popular movies
// Popular, movie sorted by rating from transaction
func (m *MovieRepository) GetPopular(reqCntxt context.Context, offset, limit int) ([]models.MovieList, error) {
	sql := `SELECT m.id, m.poster_path, m.title, AVG(t.rating) AS avg_rating, COUNT(t.id) AS rating_count,
		json_agg(DISTINCT jsonb_build_object('genre_id', g.id, 'genre_name', g.name)) AS genres
		FROM movies m
		JOIN transactions t ON m.id = t.movies_id
		JOIN genres_movies gm ON m.id = gm.movie_id
		JOIN genres g ON gm.genre_id = g.id
		WHERE t.paid_at != null AND t.rating IS NOT NULL
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
	JOIN genres_movies gm ON m.id = gm.movie_id
	JOIN genres g ON gm.genre_id = g.id `
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
	log.Println(rows)
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
	log.Println("test", Movies)
	return Movies, nil
}

// Function query get movie details
func (m *MovieRepository) GetMovieDetails(reqCntxt context.Context, movieID string) (models.MovieDetails, error) {
	sql := `SELECT m.id, m.title, m.overview, m.poster_path, m.backdrop_path, m.release_date, m.duration, d.name AS director_name, d.id AS director_id, 
	json_agg(DISTINCT jsonb_build_object('genre_id', g.id, 'genre_name', g.name)) AS genres, 
	json_agg(DISTINCT jsonb_build_object('actor_id', a.id, 'actor_name', a.name)) AS actors
  FROM movies m
  JOIN directors d ON m.director_id = d.id
  LEFT JOIN genres_movies gm ON m.id = gm.movie_id
  LEFT JOIN genres g ON gm.genre_id = g.id
  LEFT JOIN actors_movies am ON m.id = am.movie_id
  LEFT JOIN actors a ON am.actor_id = a.id
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
