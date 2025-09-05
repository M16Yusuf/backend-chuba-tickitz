package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

type MovieHandler struct {
	movRepo *repositories.MovieRepository
}

func NewMovieHandler(movRepo *repositories.MovieRepository) *MovieHandler {
	return &MovieHandler{movRepo: movRepo}
}

func (m *MovieHandler) UpcomingMovie(ctx *gin.Context) {
	// Make pagenation using query LIMIT dan OFFSET
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	movies, err := m.movRepo.GetUpcoming(ctx.Request.Context(), offset, limit)
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

	if len(movies) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"data":    []any{},
			"page":    page,
		})
		return
	}

	ctx.JSON(http.StatusOK, models.MoviesResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      200,
		},
		Data: movies,
		Page: page,
	})
}
