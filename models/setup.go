package models

import (
	// "fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var DB *gorm.DB
var DB = ConnectDatabase()

func ConnectDatabase() *gorm.DB {
	dsn := "root:123698745leo@tcp(127.0.0.1:3306)/db1?parseTime=True&loc=Local"
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to the database")

	// 测试查询工作

	// var user User

	// DB.First(&user, "id = ?", 1)
	// log.Println(user)

	DB.AutoMigrate(&User{})

	return DB
}
