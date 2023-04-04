package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"fmt"
    "strings"

    utility "golang_gin/utility"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        splitToken := strings.Split(tokenString, "Bearer ")
        tokenString = splitToken[1]
        fmt.Println(tokenString)
        if tokenString == "" {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            filename := "secret.txt"
            secret, err := utility.ReadFile(filename)
            if err != nil {
                return nil, err
            }
            return []byte(secret), nil // replace with your own secret key
        })

        if err != nil || !token.Valid {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        c.Next()
    }
}