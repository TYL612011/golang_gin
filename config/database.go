package config

import (
	"os"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	PASSWORD = os.Getenv("PASSWORD")
	USERNAME = os.Getenv("USERNAME")
	HOST     = os.Getenv("MYSQL_SERVER")
	DATABASE = os.Getenv("DATABASE")
	PORT	 = os.Getenv("PORT")
)

func urlConnect() string{
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", USERNAME, PASSWORD, HOST, PORT, DATABASE)
}

func InitConnectToMysql() (*gorm.DB, error){
	dsn := urlConnect()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	} 
	return db, nil
}

func CloseConnectToMysql(db *gorm.DB) error{
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
		return nil
	} 
	return err
}



