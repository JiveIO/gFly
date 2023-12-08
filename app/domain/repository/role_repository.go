package repository

import (
	"app/app/domain/models"
	"app/core/db"
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

	// Define query string.
	query := `SELECT * FROM roles WHERE slug=$1;`

	// Send query to database.
	err := q.DB.Get(&role, query, slug)
	if err != nil {
		// Return empty object and error.
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
		// Return empty object and error.
		return nil, err
	}

	// Return query result.
	return roles, nil
}

// AddRoleForUserID query for adding role for given user ID.
func (q *RoleRepository) AddRoleForUserID(userID uuid.UUID, slug string) error {
	// Define role variable.
	role := models.Role{}

	// Define query string.
	query := `SELECT * FROM roles WHERE slug=$1;`

	// Send query to database.
	err := q.DB.Get(&role, query, slug)
	if err != nil {
		return err
	}

	// Define query string.
	query = `INSERT INTO user_roles(id, role_id, user_id, created_at) VALUES ($1, $2, $3, $4)`

	// Send query to database.
	_, err = q.DB.Exec(
		query,
		uuid.New(), role.ID, userID, time.Now(),
	)

	if err != nil {
		// Return only error.
		return err
	}

	// This query returns nothing.
	return nil
}
