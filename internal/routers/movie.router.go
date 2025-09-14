package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func InitMovieRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	movieRouter := router.Group("/movies")
	movieRepository := repositories.NewMovieRepository(db, rdb)
	mh := handlers.NewMovieHandler(movieRepository)

	movieRouter.GET("/upcoming", mh.UpcomingMovie)
	movieRouter.GET("/popular", mh.PopularMovie)
	movieRouter.GET("", mh.FilterMovie)
	movieRouter.GET("/:movie_id", mh.GetDetailMovie)
}
