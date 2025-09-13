package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/middleware"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

func InitAdminRouter(router *gin.Engine, db *pgxpool.Pool) {
	adminRouter := router.Group("/admin")
	adminRepository := repositories.NewAdminRepository(db)
	ah := handlers.NewAdminHandler(adminRepository)

	adminRouter.GET("/movies", middleware.VerifyToken, middleware.Access("admin"), ah.GetAllMovieAdmin)
	adminRouter.DELETE("/movies/:movie_id", middleware.VerifyToken, middleware.Access("admin"), ah.DeleteMovieByID)
}
