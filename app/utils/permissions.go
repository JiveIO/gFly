package utils

import (
	"app/app/constants"
	"app/app/domain/permissions"
	"fmt"
)

// GetPermissionsByRole func for getting permissions from a role name.
func GetPermissionsByRole(role string) ([]string, error) {
	// Define permissions variable.
	var perms []string

	// Switch given role.
	switch role {
	case constants.AdminRoleName:
		// Admin permissions (all access).
		perms = []string{
			permissions.SystemViewPermission,
			permissions.SystemChangePermission,
			permissions.SystemReportPermission,
			// user_permissions
			permissions.ProfileViewPermission,
			permissions.UpdatePasswordPermission,
		}
	case constants.ModeratorRoleName:
		// Moderator permissions.
		perms = []string{
			permissions.SystemReportPermission,
			permissions.SystemViewPermission,
			// user_permissions
			permissions.ProfileViewPermission,
			permissions.UpdatePasswordPermission,
		}
	case constants.UserRoleName:
		// User permissions.
		perms = []string{
			// user_permissions
			permissions.ProfileViewPermission,
			permissions.UpdatePasswordPermission,
		}
	case constants.GuestRoleName:
		// Guest permissions.
		perms = []string{}
	default:
		// Return error message.
		return nil, fmt.Errorf("role '%v' does not exist", role)
	}

	return perms, nil
}
