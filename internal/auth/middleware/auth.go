// Package middleware provides middlewares for handling authentication
// and other cross-cutting concerns in the loyalty system.
package middleware

import (
	"net/http"
	"strings"

	"github.com/elina-chertova/loyalty-system/internal/auth/handlers"
	"github.com/elina-chertova/loyalty-system/internal/security"
	"github.com/gin-gonic/gin"
)

// JWTAuth is a middleware function for the Gin framework that handles
// JWT token validation. It checks for the presence of a JWT token either
// in the Authorization header or as a cookie named "access_token".
// If the token is valid, it allows the request to proceed; otherwise,
// it returns an unauthorized status.
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenBearer := c.GetHeader("Authorization")
		if accessTokenBearer != "" {
			extractedToken := strings.Split(accessTokenBearer, "Bearer ")
			if len(extractedToken) != 2 {
				c.JSON(
					http.StatusUnauthorized,
					handlers.Response{
						Message: "Check token",
						Status:  "Invalid Authorization Token format",
					},
				)
				c.Abort()
				return
			}

			err := security.ValidateToken(extractedToken[1])
			if err != nil {
				c.AbortWithStatusJSON(
					http.StatusUnauthorized,
					handlers.Response{
						Message: err.Error(),
						Status:  "Unauthorized",
					},
				)
				return
			}

			c.Set("token", extractedToken[1])
			c.Next()
			return
		}

		accessTokenCookie, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				handlers.Response{
					Message: err.Error(),
					Status:  "Unauthorized",
				},
			)
			c.Abort()
			return
		}

		err = security.ValidateToken(accessTokenCookie)
		if err != nil {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				handlers.Response{
					Message: err.Error(),
					Status:  "Unauthorized",
				},
			)
			return
		}

		c.Set("token", accessTokenCookie)
		c.Next()
	}
}
