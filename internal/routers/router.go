package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"

	docs "github.com/m16yusuf/backend-chuba-tickitz/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool) *gin.Engine {
	// inisialisasi engine gin
	router := gin.Default()

	// swaggo configuration
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// jika route tidak ditemukan kirim response
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, models.NoRouteResponse{
			Message: "salah...",
			Status:  "Tidak ditemukan",
		})
	})

	return router
}
