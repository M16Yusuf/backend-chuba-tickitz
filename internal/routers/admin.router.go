package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/middleware"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/redis/go-redis/v9"
)

func InitAdminRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	adminRouter := router.Group("/admin")
	adminRepository := repositories.NewAdminRepository(db, rdb)
	ah := handlers.NewAdminHandler(adminRepository)

	adminRouter.GET("/movies", middleware.VerifyToken(rdb), middleware.Access("admin"), ah.GetAllMovieAdmin)
	adminRouter.DELETE("/movies/:movie_id", middleware.VerifyToken(rdb), middleware.Access("admin"), ah.DeleteMovieByID)
	adminRouter.POST("/movies", middleware.VerifyToken(rdb), middleware.Access("admin"), ah.AddMovie)
}
