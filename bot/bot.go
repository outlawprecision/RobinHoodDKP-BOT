package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/outlawprecision/RobinHoodDKP-BOT/tree/main/db"
	"github.com/outlawprecision/RobinHoodDKP-BOT/tree/main/utils"
)

type Bot struct {
	Session *discordgo.Session
	DB      *db.DB
}

func NewBot(session *discordgo.Session, db *db.DB) *Bot {
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
	command := strings.ToLower(args[0][1:])

	switch command {
	case "checkbalance":
		b.checkBalance(s, m)
	case "addpoints":
		b.addPoints(s, m, args)
	case "removepoints":
		b.removePoints(s, m, args)
	case "watchevent":
		b.watchEvent(s, m, args)
	default:
		s.ChannelMessageSend(m.ChannelID, "Unknown command.")
	}
}

func (b *Bot) checkBalance(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Implement check balance logic
}

func (b *Bot) addPoints(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
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

func (b *Bot) removePoints(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
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
