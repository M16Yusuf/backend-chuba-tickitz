package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

func InitSeatRouter(router *gin.Engine, db *pgxpool.Pool) {
	seatRouter := router.Group("/seats")
	seatRepository := repositories.NewSeatRepository(db)
	sh := handlers.NewSeatHandler(seatRepository)

	seatRouter.GET("/:schedule_id", sh.GetBookedSeat)
}
