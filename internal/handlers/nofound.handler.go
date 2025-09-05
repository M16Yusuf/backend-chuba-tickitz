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
		Response: models.Response{
			IsSuccess: false,
			Code:      404,
		},
		Err: "Route not found ...",
	})
}
