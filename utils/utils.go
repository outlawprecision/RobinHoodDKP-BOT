package utils

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
