package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	// "strings"
	// "time"
	"log"

	config "golang_gin/config"
	utility "golang_gin/utility"
)

func AuthAccessTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, config.MESSAGE_TOKEN_MISSING)
			c.Abort()
			return
		}

		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.JSON(401, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}
		tokenString := authHeader[7:]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			secret, err := utility.ReadFile("secret_access.txt")
			if err != nil {
				return nil, err
			}
			return []byte(secret), nil // replace with your own secret key
		})
		if err != nil && !token.Valid {
			log.Println(fmt.Sprintf("ERROR - Middle Auth : %s", err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": fmt.Sprintf("%s", err)})
			return
		}

		c.Next()
	}
}
