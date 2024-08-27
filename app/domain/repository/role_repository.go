package repository

import (
	"gfly/app/domain/models"
	"github.com/gflydev/core/errors"
	"time"

	mb "github.com/gflydev/db"       // Model builder
	qb "github.com/jiveio/fluentsql" // Query builder
)

// RoleRepository struct for queries from a Role model.
// The struct is an implementation of interface IRoleRepository
type RoleRepository struct {
}

// findOne query that getting one model by a specific field condition.
func (q *RoleRepository) findOne(field string, value any) *models.Role {
	// Create an instance of User
	m := models.Role{}

	// Get model and assign into `m` struct
	err := mb.GetModel(&m, field, value)

	// Return an empty model
	if err != nil {
		return nil
	}

	return &m
}

// GetRoleBySlug query for getting role by given slug.
func (q *RoleRepository) GetRoleBySlug(slug string) *models.Role {
	return q.findOne("slug", slug)
}

// GetRolesByUserID query for getting roles by given user ID.
func (q *RoleRepository) GetRolesByUserID(userID int) []models.Role {
	// DB Model instance
	db := mb.Instance()

	// Define role variable.
	var roles []models.Role

	_, _ = db.Select("roles.*").
		Join(qb.InnerJoin, "user_roles", qb.Condition{
			Field: "roles.id",
			Opt:   qb.Eq,
			Value: qb.ValueField("user_roles.role_id"),
		}).
		Where("user_roles.user_id", qb.Eq, userID).
		OrderBy("name", qb.Asc).
		Find(&roles)

	// Return query result.
	return roles
}

// AddRoleForUserID query for adding role for given user ID.
func (q *RoleRepository) AddRoleForUserID(userID int, slug string) error {
	// Define role variable.
	var role *models.Role
	role = q.GetRoleBySlug(slug)
	if role == nil {
		return errors.New("Role not found")
	}

	// Create new user
	userRole := models.UserRole{
		ID:        userID,
		RoleID:    role.ID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	return mb.CreateModel(&userRole)
}
