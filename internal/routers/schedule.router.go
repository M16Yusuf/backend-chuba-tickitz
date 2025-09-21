package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/middleware"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func InitScheduleRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	scheduleRouter := router.Group("/schedules")
	scheduleRepository := repositories.NewScheduleRepository(db)
	sh := handlers.NewScheduleHandler(scheduleRepository)

	scheduleRouter.GET("/:movieid", middleware.VerifyToken(rdb), middleware.Access("user", "admin"), sh.GetScheduleMovie)
}
