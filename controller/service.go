package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

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
		log.Println(fmt.Sprintf("SUCCESS - Login request: User %s (%s) login success", user.Username, ip))
		c.JSON(http.StatusOK, gin.H{
			"access_token":  model.AccessToken{Token: accessToken, Exp: time.Now().Add(time.Minute * 15).Unix()},
			"refresh_token": model.RefreshToken{Token: refreshToken},
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
	fmt.Println("User information: ", user)
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
