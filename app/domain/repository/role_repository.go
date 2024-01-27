package repository

import (
	"app/app/domain/models"
	"github.com/google/uuid"
	"time"

	mb "github.com/jiveio/fluentmodel" // Model builder
	qb "github.com/jiveio/fluentsql"   // Query builder
)

// RoleRepository struct for queries from a Role model.
// The struct is an implementation of interface IRoleRepository
type RoleRepository struct {
}

// GetRoleBySlug query for getting role by given slug.
func (q *RoleRepository) GetRoleBySlug(slug string) (models.Role, error) {
	// DB Model instance
	db := mb.Instance()
	// Error variable
	var err error

	// Define role variable.
	role := models.Role{}
	err = db.Where("slug", qb.Eq, slug).
		First(&role)

	// Return query result.
	return role, err
}

// GetRolesByUserID query for getting roles by given user ID.
func (q *RoleRepository) GetRolesByUserID(userID uuid.UUID) ([]models.Role, error) {
	// DB Model instance
	db := mb.Instance()
	// Error variable
	var err error

	// Define role variable.
	var roles []models.Role

	_, err = db.Select("roles.*").
		Join(qb.InnerJoin, "user_roles", qb.Condition{
			Field: "roles.id",
			Opt:   qb.Eq,
			Value: qb.ValueField("user_roles.role_id"),
		}).
		Where("user_roles.user_id", qb.Eq, userID.String()).
		OrderBy("name", qb.Asc).
		Find(&roles)

	// Return query result.
	return roles, err
}

// AddRoleForUserID query for adding role for given user ID.
func (q *RoleRepository) AddRoleForUserID(userID uuid.UUID, slug string) error {
	// DB Model instance
	db := mb.Instance()

	// Defer a rollback in case anything fails.
	db.Begin()
	defer func(db *mb.DBModel) {
		_ = db.Rollback()
	}(db)

	// Error variable
	var err error

	// Define role variable.
	var role models.Role
	role, err = q.GetRoleBySlug(slug)
	if err != nil {
		return err
	}

	// Create new user
	userRole := models.UserRole{
		ID:        uuid.New(),
		RoleID:    role.ID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	err = db.Create(&userRole)
	if err != nil {
		return err
	}

	// DB commit
	err = db.Commit()

	// Return query result.
	return err
}
