package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

type ScheduleHandler struct {
	sr *repositories.ScheduleRepository
}

func NewScheduleHandler(sr *repositories.ScheduleRepository) *ScheduleHandler {
	return &ScheduleHandler{sr: sr}
}

// Schedule
// @Tags 				Schedules
// @Router 			/schedules/{movieid} [GET]
// @Description Get schedules movie, for a movie
// @Param				movieid	path	string 	true 	"get schedule by this id movie"
// @produce			json
// @failure 		400		{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 	{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 	{object}	models.ScheduleResponse
func (s *ScheduleHandler) GetScheduleMovie(ctx *gin.Context) {
	movieID := ctx.Param("movieid")
	schedules, err := s.sr.GetSchedule(ctx.Request.Context(), movieID)
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

	// validate if schedules movie is return empty data
	if len(schedules) == 0 {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: true,
				Code:      404,
			},
			Err: "Empty schedule list",
		})
		return
	}

	// send data schedules movie as response
	ctx.JSON(http.StatusOK, models.ScheduleResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      200,
		},
		Data: schedules,
	})
}
