package repositories

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
)

type HistoryRepository struct {
	db *pgxpool.Pool
}

func NewHistoryRepository(db *pgxpool.Pool) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (h *HistoryRepository) GetHistory(reqContxt context.Context, userID string) ([]models.History, error) {
	sql := `SELECT t.id, t.code_ticket, t.paid_at, t.total_price, t.created_at, t.rating, 
		m.title AS movie_title, sch.schedule AS schedule_time, p.method AS payment_method,
		c.name AS cinema_name, ci.name AS city_name, ARRAY_AGG(s.code) AS seat_codes
		FROM transactions t
		JOIN movies m ON t.movies_id = m.id
		JOIN schedules sch ON t.schedule_id = sch.id
		JOIN payments p ON t.payment_id = p.id
		JOIN cinemas c ON t.cinema_id = c.id
		JOIN cities ci ON t.city_id = ci.id
		JOIN order_seat os ON t.id = os.transaction_id
		JOIN seats s ON os.seat_id = s.id
		WHERE t.user_id = $1
		GROUP BY t.id, m.title, sch.schedule, p.method, c.name, ci.name
		ORDER BY t.created_at DESC`

	rows, err := h.db.Query(reqContxt, sql, userID)
	if err != nil {
		log.Println("Internal Server Error: ", err.Error())
		return []models.History{}, err
	}
	defer rows.Close()

	// 	processing data / read rows
	var histories []models.History
	for rows.Next() {
		var history models.History
		var seatRaw []byte
		err := rows.Scan(&history.Id, &history.Code, &history.PaidAt, &history.Total, &history.CreatedAt, &history.Rating, &history.Movie.Title, &history.Schedule.Schedule, &history.Payment.Method, &history.Cinema.Name, &history.City.Name, &seatRaw)
		if err != nil {
			log.Println("Scan Error, ", err.Error())
			return []models.History{}, err
		}

		// Decode raw JSON into []Genre
		if err := json.Unmarshal(seatRaw, &history.Seat); err != nil {
			log.Println("Unmarshal Error:", err)
			return nil, err
		}

		histories = append(histories, history)
	}
	return histories, nil
}
