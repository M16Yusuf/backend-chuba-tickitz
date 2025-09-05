package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/m16yusuf/backend-chuba-tickitz/pkg"
)

func Access(roles ...string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		// ambil data claim
		claims, isExist := ctx.Get("claims")
		if !isExist {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Silahkan login kembali",
			})
			return
		}
		user, ok := claims.(pkg.Claims)
		if !ok {
			// log.Println("Cannot cast claims into pkg.claims")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error",
			})
			return
		}
		if !slices.Contains(roles, user.Role) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Anda tidak punya hak akses untuk resource ini",
			})
			return
		}
		ctx.Next()
	}
}
