package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

type AdminHandler struct {
	AdRep *repositories.AdminRepository
}

func NewAdminHandler(AdRep *repositories.AdminRepository) *AdminHandler {
	return &AdminHandler{AdRep: AdRep}
}

// Get all movies (admin)
// @Tags Admin
// @Router   		/admin/movies [GET]
// @Summary 		Get all list movies
// @Description Get Get all data movies, admin role required
// @Param				page	query		int 	false 	"opsional query for pagination"
// @Security 		JWTtoken
// @produce			json
// @failure 		400		{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 	{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 	{object}	models.MoviesResponse
func (ah *AdminHandler) GetAllMovieAdmin(ctx *gin.Context) {
	// Make pagenation using query LIMIT dan OFFSET
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	// get data movies from database/repositories
	movies, err := ah.AdRep.GetAllMovies(ctx.Request.Context(), offset, limit)
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
	if len(movies) == 0 {
		ctx.JSON(http.StatusNotFound, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: true,
				Code:      404,
				Page:      page,
			},
			Err: "Empty movie list",
		})
		return
	}

	// send data movies as response
	ctx.JSON(http.StatusOK, models.MoviesResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      200,
			Page:      page,
		},
		Data: movies,
	})
}

// Delete a movie
// @Tags Admin
// @Router /admin/movies/{movie_id} [DELETE]
// @Summary 							Delete a movie
// @Description 					Delete a movie (soft delete), admin role required
// @Param				movie_id	path 	int 	true 	"movie with this id will be delete"
// @Security 		JWTtoken
// @produce			json
// @failure 		400		{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 	{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 	{object}	models.Response
func (ah *AdminHandler) DeleteMovieByID(ctx *gin.Context) {
	movieID := ctx.Param("movie_id")
	log.Println(movieID)
	if err := ah.AdRep.DeleteMovie(ctx.Request.Context(), movieID); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      400,
			},
			Err: err.Error(),
		})
		return
	}

	// send http status success delet
	ctx.JSON(http.StatusOK, models.Response{
		IsSuccess: true,
		Code:      200,
		Msg:       "Movies deleted successfully",
	})
}
