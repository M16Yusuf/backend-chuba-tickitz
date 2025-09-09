package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/m16yusuf/backend-chuba-tickitz/pkg"
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
// @Summary 		get profile user
// @Description Get details user, gt data by known id user
// @Param				user_id	path  		string	true 				 "get detail movie by id movie"
// @Param 			Authorization 		header 	string true  "Bearer token"
// @produce			json
// @failure 		400			{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 		{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 		{object}	models.UserDetailResponse
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

// Update User
// @Tags Profile
// @Router 			/users/{user_id}  [PATCH]
// @Summary 		Update registerd user
// @Description Update user and show new updated data
// @Param				user_id					path  	string			true "get user_id for select which user will update"
// @Param 			Authorization 	header 	string 			true "Bearer token"
// @Param 			body 	body 			models.User true "Data new user"
// @produce			json
// @failure 		400		{object} 	models.ErrorResponse "Bad Request"
// @failure 		500 	{object} 	models.ErrorResponse "Internal Server Error"
// @success			200 	{object}	models.ProfileResponse
func (u *UserHandler) UpdateUser(ctx *gin.Context) {
	var body models.User
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: err.Error(),
		})
		return
	}

	userID := ctx.Param("user_id")
	// hashing new password
	if body.Password != "" {
		hc := pkg.NewHashConfig()
		hc.UseRecommended()
		hash, err := hc.GenHash(body.Password)
		if err != nil {
			log.Println("Failed hashing password, \nCause: ", err)
			ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      500,
				},
				Err: err.Error(),
			})
			return
		}
		body.Password = hash
	}

	newProfile, err := u.ur.EditUser(ctx.Request.Context(), body, userID)
	if err != nil {
		log.Println("Failed function editUser, \nCause: ", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: err.Error(),
		})
		return
	}

	// return new updated user as response
	ctx.JSON(http.StatusOK, models.ProfileResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      200,
		},
		Data: newProfile,
	})
}
