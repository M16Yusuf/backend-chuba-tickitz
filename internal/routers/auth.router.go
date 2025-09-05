package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

func InitAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")
	authRepository := repositories.NewAuthRepository(db)
	authHandler := handlers.NewAuthHandler(authRepository)

	authRouter.POST("", authHandler.Login)
	authRouter.POST("/register", authHandler.Register)
}
