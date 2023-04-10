package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"

	// config "golang_gin/config"
	config "golang_gin/config"
	controller "golang_gin/controller"
	utility "golang_gin/utility"
)

func init() {
	//Load .env file (variable enviroment)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERROR - Error loading .env file")
	}

	// Settiing log
	logPath := os.Getenv("LOGPATH")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("ERROR - Failed to open log file: ", err)
	}

	loc, err := time.LoadLocation(config.TIME_ZONE)
	if err != nil {
		log.Println(fmt.Sprintf("ERROR - Read time zone with error %s", err))
	}
	log.SetOutput(file)
	log.SetFlags(0)
	log.SetPrefix(time.Now().In(loc).Format("2006-01-02 15:04:05") + " ")

	// Init secret and write to file

	secret_access, err := utility.RandSecret()
	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR - Init secret access token key with error : %s", err))
		return
	}
	utility.WriteSecretToFile(secret_access, "secret_access.txt")

	secret_refresh, err := utility.RandSecret()
	if err != nil {
		log.Fatal(fmt.Sprintf("ERROR - Init secret refresh token key with error : %s", err))
		return
	}
	utility.WriteSecretToFile(secret_refresh, "secret_refresh.txt")

	log.Println(fmt.Sprintf("SUCCESS - Init secret key for token success"))

	// Init migrate database
	err = controller.InitDB()
	if err != nil {
		log.Fatal(fmt.Sprintf("SUCCESS - Init migrate database with error : %s", err))
		return
	}
	log.Println(fmt.Sprintf("SUCCESS - Init migrate database success"))
}

func main() {
	// Main service
	router := gin.Default()
	router.POST("/login", controller.LoginHandler)
	router.POST("/register", controller.RegisterHandler)
	router.GET("/token/refresh", controller.RenewAccessToken)
	routerAuth := router.Group("/auth")
	routerAuth.Use(controller.AuthAccessTokenMiddleware())
	{
		routerAuth.GET("/protect", controller.Protected)
	}
	router.Run(":8008")
}
