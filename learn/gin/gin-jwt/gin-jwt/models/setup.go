package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	DbHost := os.Getenv("DB_HOST")
	DbPort := os.Getenv("DB_PORT")
	DbUser := os.Getenv("DB_USER")
	DbPass := os.Getenv("DB_PASS")
	DbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai password=%s", DbHost, DbPort, DbUser, DbName, DbPass)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Connect to database failed, %v\n", err)
	} else {
		log.Printf("Connect to database success, host: %s, port: %s, user: %s, dbname: %s\n", DbHost, DbPort, DbUser, DbName)
	}

	// 迁移数据表
	DB.AutoMigrate(&User{})
}
