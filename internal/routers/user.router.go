package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/handlers"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/middleware"
	"github.com/m16yusuf/backend-chuba-tickitz/internal/repositories"
)

func InitUserRouter(router *gin.Engine, db *pgxpool.Pool) {
	userRouter := router.Group("/users")
	userRepository := repositories.NewUserRepository(db)
	uh := handlers.NewUserHandler(userRepository)

	userRouter.GET("", middleware.VerifyToken, middleware.Access("user", "admin"), uh.GetUserByID)
	userRouter.PATCH("", middleware.VerifyToken, middleware.Access("user", "admin"), uh.UpdateUser)
	userRouter.PATCH("/avatar", middleware.VerifyToken, middleware.Access("user", "admin"), uh.UpdateAvatar)
	userRouter.PATCH("/password", middleware.VerifyToken, middleware.Access("user", "admin"), uh.UpdatePassword)
}
