package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"time"

	config "golang_gin/config"
	model "golang_gin/model"
	utility "golang_gin/utility"
)

func LoginHandler(c *gin.Context) {
	var user, userVerify model.User 
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusNotFound, config.MESSAGE_DATA_ERRNOTFOUND)
		return
	}

	status, _ := utility.CheckPassword(user.Password, userVerify.Password)
	if status == false {
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
		c.JSON(http.StatusOK, gin.H{
			"access_token": model.AccessToken{Token: accessToken, Exp: time.Now().Add(time.Minute * 15).Unix()},
			"refresh_token": model.RefreshToken{Token: refreshToken},
		})
	} else {
		c.JSON(http.StatusNotFound, config.MESSAGE_DATA_ERRNOTFOUND)
	}
}

func Protected(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "protected"})
}


func RegisterHandler(c *gin.Context) {
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
	if err != nil {
		c.JSON(http.StatusBadRequest, config.MESSAGE_DATA_ERRFOUND)
	}
	
	//Check password security requirement
	status := utility.ValidatePassword(user.Password)
	if !status {
		c.JSON(http.StatusNotFound, config.MESSAGE_PASSWORD_ERROR)
		return
	}
	hashPass, err := utility.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, config.MESSAGE_PASSWORD_ERROR)
	}
	user.Password = hashPass


}


