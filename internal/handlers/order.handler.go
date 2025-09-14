package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/models"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
	"github.com/m16yusuf/backend-chuba-tickitz/pkg"
)

type OrderHandler struct {
	orRep *repositories.OrderRepository
}

func NewOrderHandler(orRep *repositories.OrderRepository) *OrderHandler {
	return &OrderHandler{orRep: orRep}
}

// Create Order
// @Tags Order
// @router 	 		/order 	[POST]
// @Summary 		Create order
// @Description Create order with inputs : (user_id, schedule_id, payment_id, total_price, []seats{id, code})
// @Param 			body		body		models.CreateOrder true 		"Inputs : (user_id, schedule_id, payment_id, total_price, []seats{id, code})"
// @Security 		JWTtoken
// @produce			json
// @failure 		400		{object} 	models.BadRequestResponse "Bad Request"
// @failure 		500 	{object} 	models.InternalErrorResponse "Internal Server Error"
// @success			200 	{object}	models.Response
func (oh *OrderHandler) CreateOrder(ctx *gin.Context) {
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

	var body models.CreateOrder
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println("Failed binding data\nCause: ", err)
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      400,
			},
			Err: "Failed binding data...",
		})
		return
	}

	log.Println(userClaim)
	// execute query create order
	err := oh.orRep.CreateOrder(ctx.Request.Context(), body, userClaim.UserId)
	if err != nil {
		log.Println("Error Cause : ", err)
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Response: models.Response{
				IsSuccess: false,
				Code:      500,
			},
			Err: err.Error(),
		})
		return
	} else {
		ctx.JSON(http.StatusOK, models.Response{
			IsSuccess: true,
			Code:      200,
			Msg:       "Success make order",
		})
	}
}
