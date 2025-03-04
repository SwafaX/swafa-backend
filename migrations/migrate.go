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
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func MigrateUp() {
	DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	modelsToMigrate := []interface{}{&models.Item{}, &models.User{}, &models.Swap{}, &models.Chat{}, &models.Message{}}

	for _, model := range modelsToMigrate {
		if model == nil {
			fmt.Println("Skipping nil model")
			continue
		}
		modelName := reflect.TypeOf(model).Elem().Name()
		err := DB.AutoMigrate(model)
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
		if model == nil {
			fmt.Println("Skipping nil model")
			continue
		}
		modelName := reflect.TypeOf(model).Elem().Name()
		err := DB.Migrator().DropTable(model)
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
	} else if strings.TrimSpace(strings.ToLower(input)) == "down" {
		MigrateDown()
	} else {
		log.Fatalln("Invalid input")
	}
}
