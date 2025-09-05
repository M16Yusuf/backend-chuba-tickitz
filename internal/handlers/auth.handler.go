package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/utils"
	"github.com/m16yusuf/backend-chuba-tickitz/pkg"
)

type AuthHandler struct {
	ar *repositories.AuthRepository
}

func NewAuthHandler(ar *repositories.AuthRepository) *AuthHandler {
	return &AuthHandler{ar: ar}
}

// Login
// @tags 				login
// @router 	 		/auth 	[POST]
// @Description login using email and password and return as response with JWT token
// @Param 			body		body		 models.Auth true 		"Input email and password"
// @accept 			json
// @produce 		json
// @failure 		400 		{object} models.ErrorResponse "Bad Request"
// @failure 		500 		{object} models.ErrorResponse "Internal Server Error"
// @success 		200 		{object} models.TokenResponse
func (a *AuthHandler) Login(ctx *gin.Context) {
	var body models.Auth
	if err := ctx.ShouldBind(&body); err != nil {
		if strings.Contains(err.Error(), "required") {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      400,
				},
				Err: "Email dan Password harus diisi",
			})
			return
		}
		if strings.Contains(err.Error(), "min") {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      400,
				},
				Err: "Password minimum 8 karakter",
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "internal server error",
		})
		return
	}

	// Get userdata and validate that user
	user, err := a.ar.GetUserWithEmail(ctx.Request.Context(), body.Email)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      400,
				},
				Err: "email atau Password salah",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "internal server error",
		})
		return
	}

	// compare the password :
	// body.password => from http body / input user
	// user.Password => from query GetUserWithEmail
	hc := pkg.NewHashConfig()
	isMatched, err := hc.CompareHashAndPassword(body.Password, user.Password)
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		re := regexp.MustCompile("hash|crypto|argon2id|format")
		if re.Match([]byte(err.Error())) {
			log.Println("Error during Hashing")
		}
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "internal server error",
		})
		return
	}

	// if not match sen https status as response
	if !isMatched {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      400,
			},
			Err: "Nama atau Password salah",
		})
		return
	}

	// If match, generate jwt token and send as response
	claim := pkg.NewJWTClaims(user.Id, user.Role)
	jwtToken, err := claim.GenToken()
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "internal server error",
		})
		return
	}
	ctx.JSON(http.StatusOK, models.TokenResponse{
		Response: models.Response{
			IsSuccess: true,
			Code:      200,
		},
		Token: jwtToken,
	})
}

// Register
// @Tags					Register
// @Router			/auth/register [post]
// @Description	Register new user input email and password and return new data users
// @Param				body		body 		 models.Auth 	true		"Input email and password new user"
// @accept			json
// @produce			json
// @failure 		400 		{object} models.ErrorResponse "Bad Request"
// @failure 		500 		{object} models.ErrorResponse "Internal Server Error"
// @success			200			{object} models.ProfileResponse
func (a *AuthHandler) Register(ctx *gin.Context) {
	var body models.Auth

	// Binding data and show if there is error when binding data
	if err := ctx.ShouldBind(&body); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: "Failed binding data ...",
		})
		return
	}

	// validation register
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
		hash, err := hc.GenHash(body.Password)
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
		// if inputs is already valid format,
		// input and check if the email already registered
		user, err := a.ar.NewUser(ctx.Request.Context(), body.Email, hash)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				Response: models.Response{
					IsSuccess: false,
					Code:      400,
				},
				Err: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, models.ProfileResponse{
			Response: models.Response{
				IsSuccess: true,
				Code:      200,
			},
			Data: user,
		})
	}
}
