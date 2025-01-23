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
)

func init() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func MigrateUp() {
	initializers.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	modelsToMigrate := []interface{}{&models.Item{}, &models.User{}, &models.Swap{}}

	for _, model := range modelsToMigrate {
		err := initializers.DB.AutoMigrate(model)
		modelName := reflect.TypeOf(model).Elem().Name()

		if err != nil {
			fmt.Printf("❌ Failed to create table for %s: %v\n", modelName, err)
		} else {
			fmt.Printf("✅ Created table: %s\n", modelName)
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
			fmt.Printf("❌ Failed to drop table for %s: %v\n", modelName, err)
		} else {
			fmt.Printf("✅ Dropped table: %s\n", modelName)
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
