package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

type UserHandler struct {
	ur *repositories.UserRepository
}

func NewUserHandler(ur *repositories.UserRepository) *UserHandler {
	return &UserHandler{ur: ur}
}

// Get User
// @Tags Profile
// @Router 			/users/{user_id}  [GET]
// @Description Get details user, gt data by known id user
// @Param			user_id	path  		string	true 				 "get detail movie by id movie"
// @Param 		Authorization 		header 	string true  "Bearer token"
// @produce		json
// @failure 	400			{object} 	models.ErrorResponse "Bad Request"
// @failure 	500 		{object} 	models.ErrorResponse "Internal Server Error"
// @success		200 		{object}	models.UserDetailResponse
func (u *UserHandler) GetUserByID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	user, err := u.ur.GetDataUser(ctx.Request.Context(), userID)
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

	// send data detail user as response
	ctx.JSON(http.StatusOK, models.UserDetailResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      200,
		},
		Data: user,
	})
}
