package utils

import (
	"github.com/bwmarrin/discordgo"
)

func HasRole(s *discordgo.Session, guildID, userID, requiredRoleID string) (bool, error) {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false, err
	}

	for _, roleID := range member.Roles {
		if roleID == requiredRoleID {
			return true, nil
		}
	}
	return false, nil
}
