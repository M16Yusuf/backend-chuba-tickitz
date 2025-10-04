package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/utils"
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
// @Router 			/users/  [GET]
// @Summary 		get profile user
// @Description Get details user, gt data by known id user
// @Security 		JWTtoken
// @produce			json
// @failure 		400			{object} 	models.BadRequestResponse "Bad Request"
// @failure 		500 		{object} 	models.InternalErrorResponse "Internal Server Error"
// @success			200 		{object}	models.UserDetailResponse
func (u *UserHandler) GetUserByID(ctx *gin.Context) {
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
// @Router 			/users/  [PATCH]
// @Summary 		Update registerd user
// @Description Update user and show new updated data
// @Param 			body 	body 			models.User true "Data new user"
// @Security 		JWTtoken
// @produce			json
// @failure 		400		{object} 	models.BadRequestResponse "Bad Request"
// @failure 		500 	{object} 	models.InternalErrorResponse "Internal Server Error"
// @success			200 	{object}	models.ProfileResponse
func (u *UserHandler) UpdateUser(ctx *gin.Context) {
	// binding body to model user
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

// update avatar
// @Tags 				Profile
// @Router 			/users/avatar  [PATCH]
// @Summary 		Update avatar registerd user
// @Description Update user and show new updated data
// @Param 			avatar formData file true     "Upload good image"
// @Security 		JWTtoken
// @produce			json
// @failure 		400		{object} 	models.BadRequestResponse "Bad Request"
// @failure 		500 	{object} 	models.InternalErrorResponse "Internal Server Error"
// @success			200 	{object}	models.ProfileResponse
func (u *UserHandler) UpdateAvatar(ctx *gin.Context) {
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

	// get image from body
	var body models.UserAvatar
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("Internal server error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "internal server error",
		})
		return
	}

	// process the image
	var filename *string
	if body.Image != nil {
		UploadFileName, err := utils.FileUpload(ctx, body.Image, "profile")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      http.StatusBadRequest,
				},
				Err: err.Error(),
			})
			return
		}
		filename = &UploadFileName
	}

	// save to database
	user, err := u.ur.EditAvatarProfile(ctx.Request.Context(), filename, userID)
	if err != nil {
		log.Println("Internal server error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "Internal server error",
		})
		return
	}

	// returning updated user as response
	ctx.JSON(http.StatusOK, models.ProfileResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      200,
		},
		Data: user,
	})
}

// Update password
// @Tags 				Profile
// @Router 			/users/password  [PATCH]
// @Summary 		Update password registerd user
// @Description Update password user
// @Param				body	 body 		models.Auth 	true		"Input new password registered user"
// @Security 		JWTtoken
// @accept			json
// @produce			json
// @failure 		400			{object} 	models.BadRequestResponse "Bad Request"
// @failure 		500 		{object} 	models.InternalErrorResponse "Internal Server Error"
// @success			200			{object}	models.Response
func (u *UserHandler) UpdatePassword(ctx *gin.Context) {
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
		return
	}
	userID := userClaim.UserId

	// binding data from body for update password
	var body models.Auth
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("Failed binding data \nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "Internal server error",
		})
		return
	}

	// validate format password
	// must contain : character, digit, symbol, and 8 character
	if err := utils.RegisterValidation(body); err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      400,
			},
			Err: err.Error(),
		})
		return
	} else {
		// hash new password??
		hc := pkg.NewHashConfig()
		hc.UseRecommended()
		hashed, err := hc.GenHash(body.Password)
		if err != nil {
			log.Println("Failed hash new password ...")
			ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      500,
				},
				Err: err.Error(),
			})
			return
		}

		// insert into database new hashed password after hashed
		if err := u.ur.EditPasswordUser(ctx.Request.Context(), hashed, userID); err != nil {
			log.Println("Failed update database new password \nCause: ", err)
			ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      500,
				},
				Err: "failed update password",
			})
			return
		}
		ctx.JSON(http.StatusOK, models.Response{
			IsSuccess: true,
			Code:      200,
			Msg:       "success update password",
		})
	}
}
