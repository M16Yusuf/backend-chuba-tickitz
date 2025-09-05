package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
)

// NoRouteHandler
// @Tags NoRoute
// @Description if route not found, send 404 statusNotfound as response
// @Produce json
// @Router /{any} [get]
// @Failure 404 {object} models.ErrorResponse
func NoRouteHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, models.ErrorResponse{
		IsSuccess: false,
		Err:       "Route not found ...",
		Code:      404,
	})
}
