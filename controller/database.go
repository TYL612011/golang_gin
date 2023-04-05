package controller

import (
	config "golang_gin/config"
	model "golang_gin/model"
)

func InitDB() error {
	db, err := config.InitConnectToMysql()
	if err != nil {
		return err
	}
	defer config.CloseConnectToMysql(db)
	db.Table("users").AutoMigrate(&model.User{})
	return nil
}
