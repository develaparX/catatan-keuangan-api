package middleware

import (
	"livecode-catatan-keuangan/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	CheckToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

// untuk mengecek
// CheckToken implements AuthMiddleware.
func (a *authMiddleware) CheckToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		token := strings.Replace(header, "Bearer ", "", -1)
		claims, err := a.jwtService.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
			return
		}
		ctx.Set("userId", claims["userId"])

		ctx.Next()
	}
}

func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
