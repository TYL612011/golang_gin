package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"fmt"
	"log"
	"time"
	"os"

	controller "golang_gin/controller"
	utility "golang_gin/utility"
	config "golang_gin/config"
	
	// config "golang_gin/config"
)

func init() {
	//Load .env file (variable enviroment)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Thiết lập ghi log vào file "app.log"
	logPath := os.Getenv("LOGPATH")
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
	log.SetOutput(file)
	log.SetFlags(0)
	// Thiết lập tiền tố cho thông tin log
	log.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " ")


	//Init secret and write to file 
	filename := config.FileSecretToken
	content, err := utility.RandSecret()
	if err != nil {
		log.Fatal(fmt.Sprintf("Init secret key for token with error : %s",err))
		return
	}
	utility.WriteSecretToFile(content, filename)
	fmt.Println("1.INIT SECRET KEY SUCCESS")
	log.Println(fmt.Sprintf("Init secret key for token success"))

	// Init migrate database
	fmt.Println("start init migrate")
	err = controller.InitDB()
	if err != nil {
		log.Fatal(fmt.Sprintf("Init migrate database with error : %s",err))
		return
	}
	fmt.Println("2.INIT DATABASE SUCCESS")
	log.Println(fmt.Sprintf("Init migrate database success"))
}

func main() {
	// Service provide
	router := gin.Default()
	router.POST("/login", controller.LoginHandler)
	router.POST("/register", controller.RegisterHandler)
	routerAuth := router.Group("/auth")
	routerAuth.Use(controller.AuthMiddleware())
	{
		routerAuth.GET("/protect", controller.Protected)
	}
	router.Run(":8008")
}