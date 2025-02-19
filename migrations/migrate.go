package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/SwafaX/swafa-backend/initializers"
	"github.com/SwafaX/swafa-backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error

	config, err := initializers.LoadConfig(".")

	// https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
	dsn := fmt.Sprintf(`
		host=%s 
		user=%s 
		password=%s 
		dbname=%s 
		port=%s 
		sslmode=disable 
		TimeZone=Asia/Shanghai`,
		"localhost",
		config.DBUser,
		config.DBPass,
		config.DBName,
		"5435",
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database.")
	}
}

func MigrateUp() {
	initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	modelsToMigrate := []interface{}{&models.Item{}, &models.User{}, &models.Swap{}}

	for _, model := range modelsToMigrate {
		err := initializers.DB.AutoMigrate(model)
		modelName := reflect.TypeOf(model).Elem().Name()

		if err != nil {
			fmt.Printf("Failed to create table for %s: %v\n", modelName, err)
		} else {
			fmt.Printf("Created table: %s\n", modelName)
		}
	}

	fmt.Println("👍 Migration complete")
}

func MigrateDown() {
	modelsToDrop := []interface{}{&models.User{}, &models.Item{}, &models.Swap{}}

	for _, model := range modelsToDrop {
		err := initializers.DB.Migrator().DropTable(model)
		modelName := reflect.TypeOf(model).Elem().Name()
		if err != nil {
			fmt.Printf("Failed to drop table for %s: %v\n", modelName, err)
		} else {
			fmt.Printf("Dropped table: %s\n", modelName)
		}
	}

	fmt.Println("👍 Successfully dropped tables")
}

func main() {
	fmt.Print("What migrations type do you want to make (up/down)?: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	if strings.TrimSpace(strings.ToLower(input)) == "up" {
		MigrateUp()
		return
	} else if strings.TrimSpace(strings.ToLower(input)) == "down" {
		MigrateDown()
		return
	} else {
		log.Fatalln("Invalid input")
		return
	}
}
