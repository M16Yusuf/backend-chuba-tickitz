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

// UpComing
// @Tags 				Movies
// @Router 			/movies/upcoming [GET]
// @Summary 		Get upciming movies
// @Description Get upcoming movies, filter movies that not aired yet
// @Param				page	query		int 	false 	"opsional query for pagination"
// @produce			json
// @failure 		400		{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 	{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 	{object}	models.MoviesResponse
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

// Popular
// @Tags 				Movies
// @Router 			/movies/popular [GET]
// @Summary 		Get popular movies
// @Description Get popular movies, filter movies already rated on every transaction
// @Param				page	query		int 	false 	"opsional query for pagination"
// @produce			json
// @failure 		400		{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 	{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 	{object}	models.MoviesResponse
func (m *MovieHandler) PopularMovie(ctx *gin.Context) {
	// Make pagenation using query LIMIT dan OFFSET
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	// call function query from repository to get popular movies
	movies, err := m.movRepo.GetPopular(ctx.Request.Context(), offset, limit)
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

// Filter Search and genres
// @Tags 				Movies
// @Router 			/movies/ [GET]
// @Summary 		filter movies by genres, title and pagination
// @Description Get popular movies, filter movies by title or genres
// @Param				page		query		int 		 false 	"opsional query for pagination"
// @Param				search	query		string 	 false 	"opsional query for search title"
// @Param				genres	query		[]string false 	"opsional query for filter genres" collectionFormat(multi)
// @Produce 		json
// @produce			json
// @failure 		400			{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 		{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 		{object}	models.DetailsMovieResponse
func (m *MovieHandler) FilterMovie(ctx *gin.Context) {
	// Make pagenation using query LIMIT dan OFFSET
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}
	limit := 20
	offset := (page - 1) * limit

	// search and genre filter
	search := ctx.Query("search")
	genres := ctx.QueryArray("genres")

	// call function query from repository to get popular movies
	movies, err := m.movRepo.GetFiltered(ctx.Request.Context(), offset, limit, search, genres)
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

// Get Movie details
// @Tags Movies
// @Router			/movies/{movie_id} [GET]
// @Summary 		Get details from a movie
// @Description Get details movies, get data by known an id movie
// @Param				movie_id	path  string	true "get detail movie by id movie"
// @produce			json
// @failure 		400				{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 			{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 			{object}	models.ScheduleResponse
func (m *MovieHandler) GetDetailMovie(ctx *gin.Context) {
	movieID := ctx.Param("movie_id")
	movieDetails, err := m.movRepo.GetMovieDetails(ctx.Request.Context(), movieID)
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

	// send data details a movie as response
	ctx.JSON(http.StatusOK, models.DetailsMovieResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      200,
		},
		Data: movieDetails,
	})
}
