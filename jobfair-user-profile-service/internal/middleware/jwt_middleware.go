package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Missing Authorization header",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid Authorization format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenMalformed
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Get user_id from claims
			if userIDStr, ok := claims["user_id"].(string); ok {
				userID, err := uuid.Parse(userIDStr)
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{
						"success": false,
						"message": "Invalid user_id in token",
					})
					c.Abort()
					return
				}
				c.Set("user_id", userID)
			} else if userIDFloat, ok := claims["user_id"].(float64); ok {
				// Handle if user_id is numeric
				c.Set("user_id", uint(userIDFloat))
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "Invalid user_id in token",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
