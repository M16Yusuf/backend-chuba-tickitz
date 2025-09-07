package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/middleware"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

func InitHistoryRouter(router *gin.Engine, db *pgxpool.Pool) {
	historyRouter := router.Group("/histories")
	historyRepository := repositories.NewHistoryRepository(db)
	hh := handlers.NewHistoryHandler(historyRepository)

	historyRouter.GET("/:user_id", middleware.VerifyToken, middleware.Access("user", "admin"), hh.GetHistory)
}
