package repository

import (
	"app/app/domain/models"
	"github.com/google/uuid"
)

// =============================================
//
//	All repository interfaces
//
// =============================================
type (
	// IUserRepository an interface for any repository implementation.
	IUserRepository interface {
		GetUserByID(id uuid.UUID) (models.User, error)
		GetUserByEmail(email string) (models.User, error)
		CreateUser(u *models.User) error
		UpdateUser(u *models.User) error
		DeleteUser(u *models.User) error
		SelectUser(page, limit uint64) ([]*models.User, int, error)
	}

	// IRoleRepository an interface for any repository implementation.
	IRoleRepository interface {
		GetRoleBySlug(slug string) (models.Role, error)
		GetRolesByUserID(userID uuid.UUID) ([]models.Role, error)
		AddRoleForUserID(userID uuid.UUID, slug string) error
	}
)
