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
		fmt.Println("Error: DISCORD_TOKEN environment variable not set")
		return
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session:", err)
		return
	}

	// Initialize the database connection
	database, err := db.NewDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer database.Close()
	fmt.Printf("Database connected! database url: %s", os.Getenv("DATABASE_URL"))

	// Initialize the bot
	myBot := bot.NewBot(discord, database)

	// Start the Bot
	err = myBot.Start()
	if err != nil {
		fmt.Println("Error starting the bot:", err)
	}
	defer discord.Close()

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
