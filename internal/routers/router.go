package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/middleware"

	docs "github.com/m16yusuf/backend-chuba-tickitz/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	// inisialisasi engine gin
	router := gin.Default()
	router.Use(middleware.CORSMiddleware)

	// swaggo configuration
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// setup routing
	InitAuthRouter(router, db)
	InitMovieRouter(router, db)
	InitScheduleRouter(router, db)
	InitSeatRouter(router, db)
	InitUserRouter(router, db)

	router.NoRoute(handlers.NoRouteHandler)
	return router
}
