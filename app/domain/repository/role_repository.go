package repository

import (
	"app/app/domain/models"
	"app/core/db"
	"app/core/fluentsql"
	"github.com/google/uuid"
	"time"
)

// RoleRepository struct for queries from a Role model.
// The struct is an implementation of interface IRoleRepository
type RoleRepository struct {
	DB *db.DB
}

// GetRoleBySlug query for getting role by given slug.
func (q *RoleRepository) GetRoleBySlug(slug string) (models.Role, error) {
	// Define role variable.
	role := models.Role{}

	// Get role by slug
	err := q.DB.FluentGet(func(query fluentsql.SelectBuilder) fluentsql.SelectBuilder {
		return query.
			From(models.TableRole).
			Where(fluentsql.Eq{"slug": slug})
	}, &role)

	if err != nil {
		// Return an empty object and error.
		return role, err
	}

	// Return query result.
	return role, nil
}

// GetRolesByUserID query for getting roles by given user ID.
func (q *RoleRepository) GetRolesByUserID(userID uuid.UUID) ([]models.Role, error) {
	// Define role variable.
	var roles []models.Role

	// Define query string. Try to get high role first: Admin > Moderator > User.
	query := `SELECT r.* FROM roles r INNER JOIN user_roles ur ON ur.role_id = r.id WHERE ur.user_id=$1 ORDER BY name`

	// Send query to database.
	err := q.DB.Select(&roles, query, userID)
	if err != nil {
		// Return an empty object and error.
		return nil, err
	}

	// Return query result.
	return roles, nil
}

// AddRoleForUserID query for adding role for given user ID.
func (q *RoleRepository) AddRoleForUserID(userID uuid.UUID, slug string) error {
	// Define role variable.
	role := models.Role{}

	// Get role by slug
	err := q.DB.FluentGet(func(query fluentsql.SelectBuilder) fluentsql.SelectBuilder {
		return query.
			From(models.TableRole).
			Where(fluentsql.Eq{"slug": slug})
	}, &role)

	if err != nil {
		return err
	}

	// Insert user
	_, err = q.DB.FluentInsert(func(query fluentsql.InsertBuilder) fluentsql.InsertBuilder {
		return query.
			Columns("id", "role_id", "user_id", "created_at").
			Values(uuid.New(), role.ID, userID, time.Now())
	}, models.TableRole)

	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
