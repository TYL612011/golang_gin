package config

import (
	"github.com/gin-gonic/gin"
)

var (
	MESSAGE_DATA_BADREQUEST  = gin.H{"message": "Your data is invalid"}
	MESSAGE_DATA_ERRNOTFOUND = gin.H{"message": "Your information is incorrect"}
	MESSAGE_DATA_ERRFOUND    = gin.H{"message": "Your information is exist, create other"}
	MESSAGE_DATA_ERRPROCESS  = gin.H{"message": "Server can't proces your request"}
	MESSAGE_CONNECT_ERROR    = gin.H{"message": "Can't connect to resource"}
	MESSAGE_PASSWORD_ERROR   = gin.H{"message": "Password must meet complexity requirements (A minimum 8 characters password contains a combination of uppercase, lowercase letter, number and special character (@$!%*?&) are required.)"}
	MESSAGE_SUCCESS          = gin.H{"message": "success"}
	MESSAGE_TOKEN_EXPIRE	 = gin.H{"message": "Token expired"}
	MESSAGE_TOKEN_MISSING	 = gin.H{"message": "Authentication Header Missing"}
	MESSAGE_TOKEN_INVALID	 = gin.H{"message": "Token is invalid"}
)
