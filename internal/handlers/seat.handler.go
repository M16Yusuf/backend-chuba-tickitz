package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

type SeatHandler struct {
	seatRepo *repositories.SeatRepository
}

func NewSeatHandler(seatRepo *repositories.SeatRepository) *SeatHandler {
	return &SeatHandler{seatRepo: seatRepo}
}

// Booked seat
// @Tags 				Seats
// @Router 			/seats/{schedule_id}  [GET]
// @Description Get seat that booked  from a movie get by schedule id
// @Param				schedule_id  path		string 	true 	"get booked seat by this schedule id"
// @produce			json
// @failure 		400		{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 	{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 	{object}	models.SeatResponse
func (s *SeatHandler) GetBookedSeat(ctx *gin.Context) {
	// get schedule id
	scheduleID := ctx.Query("schedule_id")

	seats, err := s.seatRepo.GetBooked(ctx.Request.Context(), scheduleID)
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

	// validate if seats is return empty data
	if len(seats) == 0 {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: true,
				Code:      404,
			},
			Err: "Empty seat list",
		})
		return
	}

	// send data seats as response
	ctx.JSON(http.StatusOK, models.SeatResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      200,
		},
		Data: seats,
	})
}
