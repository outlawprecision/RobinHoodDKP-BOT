package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/outlawprecision/RobinHoodDKP-BOT/bot"
	"github.com/outlawprecision/RobinHoodDKP-BOT/db"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fmt.Println("Error: DISCORD_TOKEN not found in environment variables")
		return
	}

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Initialize the database
	db, err := db.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Initialize the bot
	myBot := bot.NewBot(session, db)

	// Start the bot
	err = myBot.Start()
	if err != nil {
		fmt.Println("Error starting the bot:", err)
		return
	}

	// Wait for a CTRL+C or other termination signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close the Discord session
	err = session.Close()
	if err != nil {
		fmt.Println("Error closing Discord session:", err)
	}
}
