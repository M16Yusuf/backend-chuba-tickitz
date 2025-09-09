package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/m16yusuf/backend-chuba-tickitz/pkg"
)

type HistoryHandler struct {
	hisRep *repositories.HistoryRepository
}

func NewHistoryHandler(hisRep *repositories.HistoryRepository) *HistoryHandler {
	return &HistoryHandler{hisRep: hisRep}
}

// History
// @Tags	Histories
// @Router 			/histories [GET]
// @Summary 		Get histories user
// @Description Get all list histories from a user
// @Security 		JWTtoken
// @produce			json
// @failure 		400		{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 	{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 	{object}	models.HistoiesResponse
func (h *HistoryHandler) GetHistory(ctx *gin.Context) {
	// get user_id by parsing jwt token
	claims, isExist := ctx.Get("claims")
	if !isExist {
		ctx.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      403,
			},
			Err: "Silahkan login kembali",
		})
		return
	}

	userClaim, ok := claims.(pkg.Claims)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "Internal server serror",
		})
	}

	userID := userClaim.UserId
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
