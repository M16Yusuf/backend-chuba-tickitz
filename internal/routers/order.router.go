package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/middleware"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func InitOrderRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	orderRouter := router.Group("/order")
	orderRepository := repositories.NewOrderRepository(db)
	oh := handlers.NewOrderHandler(orderRepository)

	orderRouter.POST("", middleware.VerifyToken(rdb), middleware.Access("user"), oh.CreateOrder)
}
