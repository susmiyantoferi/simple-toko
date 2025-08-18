package middleware

import (
	"fmt"
	"net/http"
	"os"
	"simple-toko/helper"
	"simple-toko/web"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			helper.ToResponseJson(ctx, http.StatusUnauthorized, "unauthorization", nil)
			ctx.Abort()
			return
		}

		tokenPart := strings.Split(authHeader, " ")
		if len(tokenPart) != 2 || tokenPart[0] != "Bearer" {
			helper.ToResponseJson(ctx, http.StatusUnauthorized, "invalid authorization format", nil)
			ctx.Abort()
			return
		}

		tokenStr := tokenPart[1]
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		claim := &web.TokenClaim{}

		token, err := jwt.ParseWithClaims(tokenStr, claim, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			helper.ToResponseJson(ctx, http.StatusUnauthorized, "invalid token", nil)
			ctx.Abort()
			return
		}

		ctx.Set("user", claim)
		ctx.Next()
	}

}

func RoleAccessMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userClaims, exist := ctx.Get("user")
		if !exist {
			helper.ToResponseJson(ctx, http.StatusUnauthorized, "user not found", nil)
			ctx.Abort()
			return
		}

		user := userClaims.(*web.TokenClaim)
		role := user.Role

		for _, v := range allowedRoles {
			if role == v {
				ctx.Next()
				return
			}
		}

		helper.ToResponseJson(ctx, http.StatusForbidden, "insufficient role", nil)
		ctx.Abort()
	}
}
