package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	// er := godotenv.Load()
	// if er != nil {
	// 	log.Fatalf("Error loading .env file: %s", er)
	// }
	// cs := os.Getenv("dsn")
	dsn := "root:Root123@tcp(127.0.0.1:3306)/golang?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf(" Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf(" Failed to configure connection pool: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println(" Database connected successfully!")
	fmt.Println("Connection pool configured successfully!")

	return DB
}
