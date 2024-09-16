package main

import (
	"database/sql"
	"log"

	"botus/internal/app"
	"botus/internal/config"
	"botus/internal/service"

	_ "github.com/lib/pq"
)

const maxRequestsPerHour = 100 // Set your limit here

func main() {

	configData, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	db, err := sql.Open("postgres", configData.DBConnectionString)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	defer db.Close()

	service.InitializeDatabase(db)

	hashService := service.NewHashService(db, maxRequestsPerHour)

	bot, err := app.NewBot(configData.BotToken, hashService)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.Start()
}
