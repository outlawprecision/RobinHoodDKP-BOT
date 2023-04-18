package utils

import (
	"fmt"
	"regexp"
	"strconv"
)

func IsAuthorized(userRoles []string, authorizedRoles []string) bool {
	// Iterate through userRoles and check if any of them match an authorized role
	for _, userRole := range userRoles {
		for _, authorizedRole := range authorizedRoles {
			if userRole == authorizedRole {
				return true
			}
		}
	}
	return false
}

func ParseDkpReward(description string) (int, error) {
	re := regexp.MustCompile(`(?i)DKP:\s*(\d+)`)
	matches := re.FindStringSubmatch(description)

	if len(matches) < 2 {
		return 0, fmt.Errorf("DKP reward not found in the event description")
	}

	dkpReward, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("Error parsing DKP reward: %v", err)
	}

	return dkpReward, nil
}
