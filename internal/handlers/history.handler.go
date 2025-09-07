package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

type HistoryHandler struct {
	hisRep *repositories.HistoryRepository
}

func NewHistoryHandler(hisRep *repositories.HistoryRepository) *HistoryHandler {
	return &HistoryHandler{hisRep: hisRep}
}

// History
// @Tags	Histories
// @Router 			/histories/{user_id} [GET]
// @Description Get all list histories from a user
// @Param				user_id	path	string 	true 	"get all list histories, by user id"
// @Param 			Authorization header string true "Bearer token"
// @produce			json
// @failure 		400		{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 	{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 	{object}	models.HistoiesResponse
func (h *HistoryHandler) GetHistory(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	histories, err := h.hisRep.GetHistory(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: err.Error(),
		})
		return
	}

	// validate if movies is return empty data
	if len(histories) == 0 {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: true,
				Code:      404,
			},
			Err: "Empty histories",
		})
		return
	}

	// send data histories as response
	ctx.JSON(http.StatusOK, models.HistoiesResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      200,
		},
		Data: histories,
	})
}
