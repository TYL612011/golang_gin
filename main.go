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
	// filename := "secret.txt"
	// content, err := utility.RandSecret()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// utility.WriteSecretToFile(content, filename)

	// Init migrate database
	// err := controller.InitDB()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println("Done")
	// hashPass, _ := utility.HashPassword("123456")
	// user1:= model.User{Username: "ThaiNHc", Password: hashPass, Email: "thaibk.nh0601@gmail.com"}
	// user2:= model.User{Username: "BacNTB", Password: hashPass, Email: "thaibk.nh0601@gmail.com"}
	db, _ := config.InitConnectToMysql()
	defer config.CloseConnectToMysql(db)
	var userVerify model.User
	db.Where("username = ?", "ThaiNH").Find(&userVerify)
	if userVerify == (model.User{}) {
		fmt.Println("Not find")
	} else {
		fmt.Println(userVerify)
	}
	

	// if userVerify != (model.User{}) {
	// 	fmt.Println("find record")
	// } else {
	// 	db.Create(&user1)
	// 	fmt.Println("Create success")
	// }
	// response2 := db.Create(user2)

	// Add some data 


	// Service provide
	// router := gin.Default()
	// router.POST("/login", controller.LoginHandler)
	// router.POST("/register", controller.RegisterHandler)
	// routerAuth := router.Group("/auth")
	// routerAuth.Use(controller.AuthMiddleware())
	// {
	// 	routerAuth.GET("/protect", controller.Protected)
	// }
	// router.Run(":8008")
}