package repository

import (
	"gfly/app/domain/models"
	"github.com/google/uuid"
)

// ====================================================================
// ======================== Repository factory ========================
// ====================================================================

// Repositories struct for collect all app repositories.
type Repositories struct {
	*UserRepository
	*RoleRepository
}

// Pool a repository pool to store all
var Pool = &Repositories{
	&UserRepository{},
	&RoleRepository{},
}

// ====================================================================
// ======================= Repository interfaces ======================
// ====================================================================

type (
	// IUserRepository an interface for any repository implementation.
	IUserRepository interface {
		GetUserByID(id uuid.UUID) *models.User
		GetUserByEmail(email string) *models.User
		CreateUser(u *models.User) error
		UpdateUser(u *models.User) error
		DeleteUser(u *models.User) error
		SelectUser(page, limit int) ([]*models.User, int, error)
	}

	// IRoleRepository an interface for any repository implementation.
	IRoleRepository interface {
		GetRoleBySlug(slug string) (models.Role, error)
		GetRolesByUserID(userID uuid.UUID) ([]models.Role, error)
		AddRoleForUserID(userID uuid.UUID, slug string) error
	}
)
