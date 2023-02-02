package models

import (
	"fmt"
	"log"
	"main/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var DB = Init()

// 数据库初始化
func Init() *gorm.DB {
	// dsn := "root:123698745leo@tcp(127.0.0.1:3306)/db1?parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True&loc=Local", config.DB_username, config.DB_password, config.DB_ip, config.DB_port, config.DB_dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	return db
}
