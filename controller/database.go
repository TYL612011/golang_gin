package controller

import (
	model "golang_gin/model"
	config "golang_gin/config"
)

func InitDB() error{
	db, err := config.InitConnectToMysql()
	if err != nil {
		return err
	}
	defer config.CloseConnectToMysql(db)
	db.Table("users").AutoMigrate(&model.User{})
	return nil
}


