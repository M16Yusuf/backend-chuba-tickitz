package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create order for user
// Insert data order
// Insert data from body (user_id, schedule_id, payment_id, total_price, []seats{id, code})
// Query tables effected : transactions, order_seat
func (o *OrderRepository) CreateOrder(reqCntxt context.Context, body models.CreateOrder, userId int) error {
	// Query Insert transaction
	sqlTransaction := `INSERT INTO transactions (user_id, schedule_id, payment_id, total_price)
		VALUES ($1, $2, $3, $4) RETURNING id`

	values := []any{userId, body.ScheduleId, body.PaymentId, body.TotalPrice}

	var tempID int
	if err := o.db.QueryRow(reqCntxt, sqlTransaction, values...).Scan(&tempID); err != nil {
		log.Println("Failed execute query\nCause: ", err)
		return err
	}

	log.Println(tempID)

	// Query insert order_seat
	sqlOrderSeat := `INSERT INTO order_seats (seat_id, transaction_id) VALUES `
	for idx, data := range body.Seats {
		sqlOrderSeat = fmt.Sprintf("%s (%d, %d)", sqlOrderSeat, data.Id, tempID)
		if idx < len(body.Seats)-1 {
			sqlOrderSeat += ", "
		}
	}
	log.Println("Final query : ", sqlOrderSeat)

	cmd, err := o.db.Exec(reqCntxt, sqlOrderSeat)
	if err != nil {
		log.Println("Failed execute query\nCause:", err)
		return err
	}
	if cmd.RowsAffected() == 0 {
		log.Println("no row effected when insert order_seat maybe failed?")
		return errors.New("no row effected when insert order_seat maybe failed?")
	}

	// if success/no error return error is nil
	return nil
}
