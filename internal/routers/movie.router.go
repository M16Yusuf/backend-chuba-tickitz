package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

func InitMovieRouter(router *gin.Engine, db *pgxpool.Pool) {
	movieRouter := router.Group("/movies")
	movieRepository := repositories.NewMovieRepository(db)
	mh := handlers.NewMovieHandler(movieRepository)

	movieRouter.GET("/upcoming", mh.UpcomingMovie)
}
