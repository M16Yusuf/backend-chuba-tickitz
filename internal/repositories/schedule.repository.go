package repositories

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
)

type ScheduleRepository struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// func to get a movie schedule
func (s *ScheduleRepository) GetSchedule(reqContxt context.Context, movieid string) ([]models.Schedule, error) {

	sql := `SELECT s.id AS schedule_id, m.id AS movie_id, m.title, s.schedule
	FROM schedules s
  JOIN movies m ON m.id = s.movie_id
  WHERE m.id = $1
  ORDER BY s.schedule ASC;`

	log.Println("movie id nih :", movieid)

	rows, err := s.db.Query(reqContxt, sql, movieid)
	if err != nil {
		log.Println("Internal Server Error: ", err.Error())
		return []models.Schedule{}, err
	}
	defer rows.Close()

	log.Println("test", rows)

	// process pgx.rows into slice of model.Schedule
	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		if err := rows.Scan(&schedule.Id, &schedule.MovieId, &schedule.MovieName, &schedule.Schedule); err != nil {
			log.Println("Scan Error, ", err.Error())
			return []models.Schedule{}, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}
