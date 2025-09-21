package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/middleware"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func InitSeatRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	seatRouter := router.Group("/seats")
	seatRepository := repositories.NewSeatRepository(db)
	sh := handlers.NewSeatHandler(seatRepository)

	seatRouter.GET("/:schedule_id", middleware.VerifyToken(rdb), middleware.Access("user", "admin"), sh.GetBookedSeat)
}
