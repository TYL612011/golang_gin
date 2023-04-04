package main

import (
	// "github.com/gin-gonic/gin"
	"fmt"

	// controller "golang_gin/controller"
	// utility "golang_gin/utility"
	model "golang_gin/model"
	config "golang_gin/config"
)

func main() {

	//Init secret and write to file 
	filename := "secret.txt"
	content, err := utility.RandSecret()
	if err != nil {
		fmt.Println(err)
		return
	}
	utility.WriteSecretToFile(content, filename)
	fmt.Println("1.INIT SECRET KEY SUCCESS")

	// Init migrate database
	err := controller.InitDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("1.INIT DATABASE SUCCESS")

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