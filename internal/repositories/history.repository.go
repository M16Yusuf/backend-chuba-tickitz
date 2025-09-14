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

func (h *HistoryRepository) GetHistory(reqContxt context.Context, userID int) ([]models.History, error) {

	sql := `SELECT t.id, t.code_ticket, t.paid_at, t.total_price, t.created_at, t.rating, 
		m.title AS movie_title, sch.schedule AS schedule_time, p.method AS payment_method,
    c.name AS cinema_name, ci.name AS city_name, 
		json_agg(DISTINCT jsonb_build_object('seat_id', s.id, 'seat_code', s.code)) AS seat_codes
		FROM transactions t
		LEFT JOIN schedules sch ON t.schedule_id = sch.id
		LEFT JOIN movies m ON sch.movie_id = m.id
		LEFT JOIN cinemas c ON sch.cinema_id = c.id
		LEFT JOIN cities ci ON sch.city_id = ci.id
		LEFT JOIN payments p ON t.payment_id = p.id
		LEFT JOIN order_seat os ON t.id = os.transaction_id
    LEFT JOIN seats s ON os.seat_id = s.id
		WHERE t.user_id = $1
		GROUP BY t.id, m.title, sch.schedule, p.method, c.name, ci.name
		ORDER BY t.created_at DESC;`

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
