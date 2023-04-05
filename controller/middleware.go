package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"fmt"
    "strings"

    utility "golang_gin/utility"
    config "golang_gin/config"
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
            filename := config.FileSecretToken
            fmt.Println("This is filename: ", filename)
            secret, err := utility.ReadFile(filename)
            if err != nil {
                return nil, err
            }
            return []byte(secret), nil // replace with your own secret key
        })
        fmt.Println("This is token information: ", token)
        if err != nil || !token.Valid {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }
        c.Next()
    }
}