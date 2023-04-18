package main

import (
	"fmt"
	"os"

	"github.com/outlawprecision/RobinHoodDKP-BOT/bot"
	"github.com/outlawprecision/RobinHoodDKP-BOT/db"
)

func main() {
	// Replace these with the appropriate values for your project
	token := os.Getenv("DISCORD_TOKEN")
	//dbHost := os.Getenv("DB_HOST")
	//dbUser := os.Getenv("DB_USER")
	//dbPassword := os.Getenv("DB_PASSWORD")
	//dbName := os.Getenv("DB_NAME")
	dbURL := os.Getenv("DATABASE_URL")

	// Connect to the PostgreSQL database
	database, err := db.NewDB(dbURL)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer database.Close()

	// Create a new Discord bot instance
	dkpBot, err := bot.NewBot(token, database)
	if err != nil {
		fmt.Println("Error creating the bot:", err)
		return
	}

	// Start the bot and listen for incoming events
	err = dkpBot.Start()
	if err != nil {
		fmt.Println("Error starting the bot:", err)
		return
	}
}
