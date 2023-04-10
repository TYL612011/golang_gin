package controller

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	// "strings"

	config "golang_gin/config"
	model "golang_gin/model"
	utility "golang_gin/utility"
)

func LoginHandler(c *gin.Context) {
	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Request.RemoteAddr
	}

	var user, userVerify model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, config.MESSAGE_DATA_BADREQUEST)
		return
	}
	//check user exist
	db, err := config.InitConnectToMysql()
	if err != nil {
		c.JSON(http.StatusInternalServerError, config.MESSAGE_CONNECT_ERROR)
		return
	}
	defer config.CloseConnectToMysql(db)
	err = db.Where("Username = ?", user.Username).First(&userVerify).Error
	if err != nil {
		log.Println(fmt.Sprintf("ERROR - Login request: User %s (%s) login failed - User doesn't exist", user.Username, ip))
		c.JSON(http.StatusNotFound, config.MESSAGE_DATA_ERRNOTFOUND)
		return
	}

	status, _ := utility.CheckPassword(user.Password, userVerify.Password)
	if status == true {
		accessToken, err := utility.GenerateAccessToken(user.Username)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		refreshToken, err := utility.GenerateRefreshToken(user.Username)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		loc, err := time.LoadLocation(config.TIME_ZONE)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, config.MESSAGE_DATA_ERRPROCESS)
		}

		log.Println(fmt.Sprintf("SUCCESS - Login request: User %s (%s) login success", user.Username, ip))
		c.JSON(http.StatusOK, gin.H{
			"access_token":  model.Token{Value: accessToken, Exp: time.Now().In(loc).Add(time.Minute * 2).Format("2006-01-02 15:04:05")},
			"refresh_token": model.Token{Value: refreshToken, Exp: time.Now().In(loc).Add(time.Hour * 24 * 7).Format("2006-01-02 15:04:05")},
		})
	} else {
		log.Println(fmt.Sprintf("ERROR - Login request: User %s (%s) login failed with incorrect password information", user.Username, ip))
		c.JSON(http.StatusNotFound, config.MESSAGE_DATA_ERRNOTFOUND)
	}
}

func Protected(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "protected"})
}

func RegisterHandler(c *gin.Context) {
	//Get request ip client
	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = c.Request.RemoteAddr
	}

	var user, userVerify model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//Check user exist
	db, err := config.InitConnectToMysql()
	if err != nil {
		c.JSON(http.StatusInternalServerError, config.MESSAGE_CONNECT_ERROR)
	}
	defer config.CloseConnectToMysql(db)
	err = db.Where("username = ?", user.Username).First(&userVerify).Error
	if err == nil {
		log.Println(fmt.Sprintf("ERROR - Request registered user : User %s (ip: %s) was failed because user exist", user.Username, ip))
		c.JSON(http.StatusBadRequest, config.MESSAGE_DATA_ERRFOUND)
		return
	}

	//Check password security requirement
	status := utility.ValidatePassword(string(user.Password))
	if !status {
		c.JSON(http.StatusNotFound, config.MESSAGE_PASSWORD_ERROR)
		return
	}
	hashPass, err := utility.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, config.MESSAGE_PASSWORD_ERROR)
	}
	user.Password = hashPass
	err = db.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, config.MESSAGE_DATA_ERRPROCESS)
		return
	}
	log.Println(fmt.Sprintf("SUCCESS - Request registered user : User %s (ip: %s) was registered", user.Username, ip))
	c.JSON(http.StatusOK, config.MESSAGE_SUCCESS)
	return
}

func RenewAccessToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, config.MESSAGE_TOKEN_MISSING)
		c.Abort()
		return
	}

	// Check if the token starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid authorization format"})
		c.Abort()
		return
	}

	// Extract the token from the Authorization header
	tokenString := authHeader[7:]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			c.JSON(http.StatusBadRequest, config.MESSAGE_TOKEN_MISSING)
			c.Abort()
		}

		secret, err := utility.ReadFile("secret_refresh.txt")
		if err != nil {
			return nil, err
		}
		return []byte(secret), nil // replace with your own secret key
	})

	if err != nil && !token.Valid {

		c.JSON(http.StatusUnauthorized, config.MESSAGE_TOKEN_INVALID)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.JSON(http.StatusBadRequest, config.MESSAGE_TOKEN_INVALID)
		return
	}

	accessToken, err := utility.GenerateAccessToken(claims["username"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, config.MESSAGE_DATA_ERRPROCESS)
		return
	}

	loc, err := time.LoadLocation(config.TIME_ZONE)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, config.MESSAGE_DATA_ERRPROCESS)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": model.Token{Value: accessToken, Exp: time.Now().In(loc).Add(time.Minute * 5).Format("2006-01-02 15:04:05")},
	})
	return
}
