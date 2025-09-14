package repositories

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
)

type SeatRepository struct {
	db *pgxpool.Pool
}

func NewSeatRepository(db *pgxpool.Pool) *SeatRepository {
	return &SeatRepository{db: db}
}

// get booked seat of a movie by a schedule id
func (s *SeatRepository) GetBooked(reqContxt context.Context, idSchedule string) ([]models.BookedSeatBySchedule, error) {
	log.Println(idSchedule)
	sql := `SELECT sch.id AS schedule_id, m.title, sch.schedule, s.code AS seat_code, t.id AS transaction_id, t.paid_at
		FROM schedules sch
		JOIN movies m ON sch.movie_id = m.id
		JOIN transactions t ON sch.id = t.schedule_id
		JOIN order_seats os ON t.id = os.transaction_id
		JOIN seats s ON os.seat_id = s.id
		WHERE sch.id = $1
		ORDER BY s.code ASC;`

	rows, err := s.db.Query(reqContxt, sql, idSchedule)
	log.Println(rows)
	if err != nil {
		log.Println("internal server error : ", err.Error())
		return []models.BookedSeatBySchedule{}, err
	}
	defer rows.Close()

	// processing data / read rows
	var seats []models.BookedSeatBySchedule
	for rows.Next() {
		var seat models.BookedSeatBySchedule
		if err := rows.Scan(&seat.ScheduleId, &seat.Title, &seat.Schedule, &seat.Code, &seat.TransactionId, &seat.PaidAt); err != nil {
			log.Println("Scan Error, ", err.Error())
			return []models.BookedSeatBySchedule{}, err
		}
		seats = append(seats, seat)
	}
	return seats, nil
}
