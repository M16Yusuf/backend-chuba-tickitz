package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/m16yusuf/backend-chuba-tickitz/pkg"
)

type AuthHandler struct {
	ar *repositories.AuthRepository
}

func NewAuthHandler(ar *repositories.AuthRepository) *AuthHandler {
	return &AuthHandler{ar: ar}
}

// Login
// @tags 			login
// @router 		/auth 	[POST]
// @Param 		body		body		 models.Auth true "Input email and password"
// @accept 		json
// @produce 	json
// @failure 	400 		{object} models.ErrorResponse "Bad Request"
// @failure 	500 		{object} models.ErrorResponse "Internal Server Error"
// @success 	200 		{object} models.TokenResponse
func (a *AuthHandler) Login(ctx *gin.Context) {
	var body models.Auth
	if err := ctx.ShouldBind(&body); err != nil {
		if strings.Contains(err.Error(), "required") {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				IsSuccess: false,
				Err:       "Email dan Password harus diisi",
				Code:      400,
			})
			return
		}
		if strings.Contains(err.Error(), "min") {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				IsSuccess: false,
				Err:       "Password minimum 8 karakter",
				Code:      400,
			})
			return
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Err:       "internal server error",
			Code:      500,
		})
		return
	}

	// Get userdata and validate that user
	user, err := a.ar.GetUserWithEmail(ctx.Request.Context(), body.Email)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
				IsSuccess: false,
				Err:       "email atau Password salah",
				Code:      400,
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Err:       "internal server error",
			Code:      500,
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
			IsSuccess: false,
			Err:       "internal server error",
			Code:      500,
		})
		return
	}

	// if not match sen https status as response
	if !isMatched {
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			IsSuccess: false,
			Err:       "Nama atau Password salah",
			Code:      400,
		})
		return
	}

	// If match, generate jwt token and send as response
	claim := pkg.NewJWTClaims(user.Id, user.Role)
	jwtToken, err := claim.GenToken()
	if err != nil {
		log.Println("Internal Server Error.\nCause: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			IsSuccess: false,
			Err:       "internal server error",
			Code:      500,
		})
		return
	}
	ctx.JSON(http.StatusOK, models.TokenResponse{
		IsSuccess: true,
		Token:     jwtToken,
	})
}
