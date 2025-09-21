package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")
	authRepository := repositories.NewAuthRepository(db, rdb)
	authHandler := handlers.NewAuthHandler(authRepository)

	authRouter.POST("", authHandler.Login)
	authRouter.POST("/register", authHandler.Register)
	authRouter.DELETE("", authHandler.Logout)
}
