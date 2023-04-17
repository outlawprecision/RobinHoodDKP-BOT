package bot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/outlawprecision/RobinHoodDKP-BOT/db"
	"github.com/outlawprecision/RobinHoodDKP-BOT/utils"
)

type Bot struct {
	Session *discordgo.Session
	DB      *db.Database
}

func NewBot(session *discordgo.Session, db *db.Database) *Bot {
	return &Bot{
		Session: session,
		DB:      db,
	}
}

func (b *Bot) Start() error {
	b.Session.AddHandler(b.onMessageCreate)

	err := b.Session.Open()
	if err != nil {
		return fmt.Errorf("error opening connection to Discord: %v", err)
	}
	return nil
}

func (b *Bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	args := strings.Fields(m.Content)
	command := args[0][1:]

	switch command {
	case "checkbalance":
		b.checkBalance(m)
	case "addpoints":
		b.addPoints(s, m, args)
	case "removepoints":
		b.removePoints(s, m, args)
	case "watchevent":
		b.watchEvent(s, m, args)
	}
}

func (b *Bot) checkBalance(m *discordgo.MessageCreate) {
	// Convert the user ID to int64
	userID, err := strconv.ParseInt(m.Author.ID, 10, 64)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, "Error: Invalid user ID")
		return
	}

	balance, err := b.DB.CheckBalance(userID)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error checking balance: %v", err))
		return
	}

	// Check if the user exists in the database
	if balance == 0 {
		err = b.DB.CreateUser(userID)
		if err != nil {
			b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error creating user: %v", err))
			return
		}
	}

	b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Your balance is: %d", balance))
}

func (b *Bot) addPoints(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !addpoints @User <amount>")
		return
	}

	// Replace this with the actual role ID for the role that you want to allow access to the addPoints command
	requiredRoleID := "1037734711201644614"

	hasRole, err := utils.HasRole(s, m.GuildID, m.Author.ID, requiredRoleID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking user role.")
		return
	}

	if !hasRole {
		s.ChannelMessageSend(m.ChannelID, "You don't have the required role to execute this command.")
		return
	}

	targetUserID, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid user.")
		return
	}
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid amount.")
		return
	}

	// Ensure the target user exists in the database
	err = b.DB.CreateUser(targetUserID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error creating target user in the database.")
		return
	}

	// Get the target user's current DP balance
	currentBalance, err := b.DB.CheckBalance(targetUserID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving target user balance.")
		return
	}

	// Update the target user's DP balance
	newBalance := currentBalance + amount
	err = b.DB.AddPoints(targetUserID, newBalance)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error updating target user balance.")
		return
	}

	// Send a confirmation message to the channel
	response := fmt.Sprintf("Successfully added %d DP to <@%s>. Their new balance is: %d", amount, targetUserID, newBalance)
	s.ChannelMessageSend(m.ChannelID, response)
}

func (b *Bot) removePoints(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Usage: !removepoints @User <amount>")
		return
	}

	// Replace this with the actual role ID for the role that you want to allow access to the removePoints command
	requiredRoleID := "1037734711201644614"

	hasRole, err := utils.HasRole(s, m.GuildID, m.Author.ID, requiredRoleID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking user role.")
		return
	}

	if !hasRole {
		s.ChannelMessageSend(m.ChannelID, "You don't have the required role to execute this command.")
		return
	}

	targetUserID, err := strconv.ParseInt(args[1], 10, 64)
	amount, err := strconv.Atoi(args[2])
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Invalid amount.")
		return
	}

	// Ensure the target user exists in the database
	err = b.DB.CreateUser(targetUserID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error creating target user in the database.")
		return
	}

	// Get the target user's current DP balance
	currentBalance, err := b.DB.CheckBalance(targetUserID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error retrieving target user balance.")
		return
	}

	// Update the target user's DP balance
	newBalance := currentBalance - amount
	if newBalance < 0 {
		s.ChannelMessageSend(m.ChannelID, "Cannot remove more points than the user has.")
		return
	}

	err = b.DB.AddPoints(targetUserID, newBalance)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error updating target user balance.")
		return
	}

	// Send a confirmation message to the channel
	response := fmt.Sprintf("Successfully removed %d DP from <@%s>. Their new balance is: %d", amount, targetUserID, newBalance)
	s.ChannelMessageSend(m.ChannelID, response)
}

func (b *Bot) watchEvent(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	hasRole, err := utils.HasRole(s, m.GuildID, m.Author.ID, "your_required_role_id")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Error checking user role.")
		return
	}

	if !hasRole {
		s.ChannelMessageSend(m.ChannelID, "You don't have the required role to execute this command.")
		return
	}

	// Continue with the command execution

}
