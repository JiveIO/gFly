package repository

import (
	"gfly/app/domain/models"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"time"

	mb "github.com/gflydev/db"          // Model builder
	qb "github.com/jivegroup/fluentsql" // Query builder
)

// ====================================================================
// ==================== Role Repository Interface =====================
// ====================================================================

// IRoleRepository an interface for any repository implementation.
type IRoleRepository interface {
	GetRolesByUserID(userID int) []models.Role
	AddRoleForUserID(userID int, slug string) error
}

// ====================================================================
// ==================== Role Repository Implement =====================
// ====================================================================

// RoleRepository struct for queries from a Role model.
// The struct is an implementation of interface IRoleRepository
type RoleRepository struct {
}

// GetRolesByUserID query for getting roles by given user ID.
func (q *RoleRepository) GetRolesByUserID(userID int) []models.Role {
	// Define role variable.
	var roles []models.Role

	_, err := mb.Instance().Select(models.TableRole+".*").
		Join(qb.InnerJoin, models.TableUserRole, qb.Condition{
			Field: models.TableRole + ".id",
			Opt:   qb.Eq,
			Value: qb.ValueField(models.TableUserRole + ".role_id"),
		}).
		Where(models.TableUserRole+".user_id", qb.Eq, userID).
		OrderBy("name", qb.Asc).
		Find(&roles)

	if err != nil {
		log.Error(err)
	}

	// Return query result.
	return roles
}

// AddRoleForUserID query for adding role for given user ID.
func (q *RoleRepository) AddRoleForUserID(userID int, slug string) error {
	// Get role by slug
	role, err := mb.GetModel[models.Role](qb.Condition{
		Field: "slug",
		Opt:   qb.Eq,
		Value: slug,
	})
	if err != nil || role == nil {
		log.Error(err)

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
