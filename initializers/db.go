package initializers

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error

	// https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
	dsn := fmt.Sprintf(`
		host=%s 
		user=%s 
		password=%s 
		dbname=%s 
		port=%s 
		sslmode=disable 
		TimeZone=Asia/Shanghai`,
		config.DBHost,
		config.DBUser,
		config.DBPass,
		config.DBName,
		config.DBPort,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database.")
	}

	fmt.Println("✔ Successfully connected to Databse.")
}
