package bot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/outlawprecision/RobinHoodDKP-BOT/db"
	"github.com/outlawprecision/RobinHoodDKP-BOT/utils"
)

type Bot struct {
	Session *discordgo.Session
	DB      *db.DB
	Config  struct {
		AdminRoleIDs []string
		ServerID     string
	}
}

type GuildEvent struct {
	ID          string `json:"id"`
	GuildID     string `json:"guild_id"`
	ChannelID   string `json:"channel_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	StartTime   string `json:"scheduled_start_time"`
	EndTime     string `json:"scheduled_end_time"`
}

func NewBot(token string, db *db.DB) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	adminRoleIDs := []string{"1037734711201644614", "<RoleID2>", "<RoleID3>"} // Replace with the actual authorized role IDs
	serverID := ""

	bot := &Bot{
		Session: session,
		DB:      db,
		Config: struct {
			AdminRoleIDs []string
			ServerID     string
		}{
			AdminRoleIDs: adminRoleIDs,
			ServerID:     serverID,
		},
	}

	return bot, nil
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
	case "enrolldkp":
		b.enrollDkp(m)
	case "checkbalance":
		b.checkBalance(m, args)
	case "addpoints":
		b.addPoints(m, args)
	case "removepoints":
		b.removePoints(m, args)
	case "dkpevent":
		b.handleDkpEvent(s, m, args)
	}
}

func (b *Bot) enrollDkp(m *discordgo.MessageCreate) {
	userID, err := strconv.ParseInt(m.Author.ID, 10, 64)
	err = b.DB.CreateUser(userID)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error creating user: %v", err))
		return
	}
	b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You have been enrolled!"))
}

func (b *Bot) checkBalance(m *discordgo.MessageCreate, args []string) {
	userID := m.Author.ID

	// Convert the author ID to int64
	authorID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, "Error: Invalid user ID")
		return
	}

	// Check if the user exists in the database
	userExists, err := b.DB.UserExists(authorID)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, "Error: Cannot check user existence")
		return
	}

	if !userExists {
		// Create the user if they don't exist
		err = b.DB.CreateUser(authorID)
		if err != nil {
			b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error creating user: %v", err))
			return
		}
	}

	// Get the user's balance
	balance, err := b.DB.CheckBalance(authorID)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error checking balance: %v", err))
		return
	}

	b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s, your balance is %d DKP", m.Author.Username, balance))
}

func (b *Bot) addPoints(m *discordgo.MessageCreate, args []string) {
	// Ensure we have enough arguments
	if len(args) < 3 {
		b.Session.ChannelMessageSend(m.ChannelID, "Usage: !addpoints <@user> <points>")
		return
	}
	// Check the user's role for authorization
	if !utils.IsAuthorized(m.Member.Roles, b.Config.AdminRoleIDs) {
		b.Session.ChannelMessageSend(m.ChannelID, "You are not authorized to use this command.")
		return
	}

	// Parse the target user ID from the message
	targetUserID := strings.TrimPrefix(strings.TrimSuffix(args[1], ">"), "<@")
	userID, err := strconv.ParseInt(targetUserID, 10, 64)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, "Error: Invalid target user ID")
		fmt.Printf("Target userID: %s", targetUserID)
		return
	}

	// Get the target user
	targetUser, err := b.Session.User(targetUserID)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error getting target user: %v", err))
		return
	}

	// Parse the points
	points, err := strconv.Atoi(args[2])
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, "Error: Invalid points value")
		return
	}
	// Check if the target user exists in the database
	userExists, err := b.DB.UserExists(userID)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, "Error: Cannot check user existence")
		return
	}

	if !userExists {
		// Create the target user if they don't exist
		err = b.DB.CreateUser(userID)
		if err != nil {
			b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error creating user: %v", err))
			return
		}
	}

	// Add points to the target user
	err = b.DB.AddPoints(userID, points)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error adding points: %v", err))
		return
	}

	b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Added %d points to %s", points, targetUser.Username))
}

func (b *Bot) removePoints(m *discordgo.MessageCreate, args []string) {
	// Ensure we have enough arguments
	if len(args) < 3 {
		b.Session.ChannelMessageSend(m.ChannelID, "Usage: !removepoints <@user> <points>")
		return
	}

	// Check the user's role for authorization
	if !utils.IsAuthorized(m.Member.Roles, b.Config.AdminRoleIDs) {
		b.Session.ChannelMessageSend(m.ChannelID, "You are not authorized to use this command.")
		return
	}

	// Parse the target user ID from the message
	targetUserID := strings.TrimPrefix(strings.TrimSuffix(args[1], ">"), "<@")

	// Convert the target user ID to int64
	userID, err := strconv.ParseInt(targetUserID, 10, 64)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, "Error: Invalid target user ID")
		return
	}

	// Get the target user
	targetUser, err := b.Session.User(targetUserID)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error getting target user: %v", err))
		return
	}

	// Parse the points
	points, err := strconv.Atoi(args[2])
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, "Error: Invalid points value")
		return
	}

	// Remove points from the target user
	err = b.DB.RemovePoints(userID, points)
	if err != nil {
		b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error removing points: %v", err))
		return
	}

	b.Session.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Removed %d points from %s", points, targetUser.Username))
}

func (b *Bot) FetchEventDetails(eventLink string) (*GuildEvent, error) {
	// Parse event link to get server ID (Guild ID), channel ID, and event ID
	eventLinkPattern := `https://discord.com/events/(\d+)/(\d+)`
	re := regexp.MustCompile(eventLinkPattern)
	matches := re.FindStringSubmatch(eventLink)
	if len(matches) != 3 {
		return nil, fmt.Errorf("Invalid event link")
	}

	guildID, eventID := matches[1], matches[2]

	// Build the Discord API URL to fetch the event details
	apiURL := fmt.Sprintf("https://discord.com/api/v10/guilds/%s/events/%s", guildID, eventID)

	// Create an HTTP request with the Discord API URL and the bot token
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bot "+b.Session.Token)
	req.Header.Set("Content-Type", "application/json")

	// Send the request and parse the response to get the event details
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error fetching event details: %s", resp.Status)
	}

	var event GuildEvent
	err = json.NewDecoder(resp.Body).Decode(&event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (b *Bot) handleDkpEvent(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	//Check if user is Authorized to use this command
	if !utils.IsAuthorized(m.Member.Roles, b.Config.AdminRoleIDs) {
		s.ChannelMessageSend(m.ChannelID, "You do not have permission to use this command.")
		return
	}

	if len(args) < 2 {
		s.ChannelMessageSend(m.ChannelID, "Please provide a link to the Discord event.")
		return
	}

	eventLink := args[1]
	eventDetails, err := b.FetchEventDetails(eventLink)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error fetching event details: %v", err))
		return
	}

	dkpReward, err := utils.ParseDkpReward(eventDetails.Description)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error parsing DKP reward: %v", err))
		return
	}

	event := &db.Event{
		EventID:     eventDetails.ID,
		Name:        eventDetails.Name,
		Description: eventDetails.Description,
		StartTime:   eventDetails.StartTime,
		EndTime:     eventDetails.EndTime,
		DKP_Reward:  dkpReward,
	}

	err = b.DB.CreateEvent(event)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error saving event to the database: %v", err))
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Event '%s' saved successfully.", eventDetails.Name))
}
