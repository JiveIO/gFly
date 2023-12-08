package repository

import (
	"app/app/domain/models"
	"app/core/db"
	"app/core/fluentsql"
	"github.com/google/uuid"
)

// UserRepository struct for queries from a User model.
// The struct is an implementation of interface IUserRepository
type UserRepository struct {
	DB *db.DB
}

// GetUserByID query for getting one User by given ID.
func (q *UserRepository) GetUserByID(id uuid.UUID) (models.User, error) {
	// Define User variable.
	user := models.User{}

	err := q.DB.FluentGet(func(query fluentsql.SelectBuilder) fluentsql.SelectBuilder {
		return query.
			From("users").
			Where(fluentsql.Eq{"id": id})
	}, &user)

	// Return empty object and error.
	if err != nil {
		return user, err
	}

	// Return query result.
	return user, nil
}

// GetUserByEmail query for getting one User by given Email.
func (q *UserRepository) GetUserByEmail(email string) (models.User, error) {
	// Define User variable.
	user := models.User{}

	err := q.DB.FluentGet(func(query fluentsql.SelectBuilder) fluentsql.SelectBuilder {
		return query.
			From("users").
			Where(fluentsql.Eq{"email": email})
	}, &user)

	if err != nil {
		// Return empty object and error.
		return user, err
	}

	// Return query result.
	return user, nil
}

// CreateUser a query for creating a new user by given user data.
func (q *UserRepository) CreateUser(u *models.User) error {
	// Define query string.
	_, err := q.DB.FluentInsert(func(query fluentsql.InsertBuilder) fluentsql.InsertBuilder {
		return query.
			Columns("id", "email", "password_hash", "fullname", "phone", "token", "user_status", "created_at", "updated_at").
			Values(u.ID, u.Email, u.PasswordHash, u.Fullname, u.Phone, u.Token, u.UserStatus, u.CreatedAt, u.UpdatedAt)
	}, "users")

	if err != nil {
		return err
	}

	// This query returns nothing.
	return nil
}

// UpdateUser a query for updating a user by given user data.
func (q *UserRepository) UpdateUser(u *models.User) error {
	// Define query string.
	_, err := q.DB.FluentUpdate(func(query fluentsql.UpdateBuilder) fluentsql.UpdateBuilder {
		return query.
			Set("email", u.Email).
			Set("password_hash", u.PasswordHash).
			Set("fullname", u.Fullname).
			Set("phone", u.Phone).
			Set("token", u.Token).
			Set("user_status", u.UserStatus).
			Set("updated_at", u.UpdatedAt).
			Where(fluentsql.Eq{"id": u.ID})
	}, "users")

	if err != nil {
		return err
	}

	// This query returns nothing.
	return nil
}

// DeleteUser a query for updating a user by given user data.
func (q *UserRepository) DeleteUser(u *models.User) error {
	// Define query string.
	_, err := q.DB.FluentDelete(func(query fluentsql.DeleteBuilder) fluentsql.DeleteBuilder {
		return query.
			Where(fluentsql.Eq{"id": u.ID})
	}, "users")

	if err != nil {
		return err
	}

	// This query returns nothing.
	return nil
}

// SelectUser a query for updating a user by given user data.
func (q *UserRepository) SelectUser(page, limit uint64) ([]*models.User, int, error) {
	// Define User variable.
	var users []*models.User
	var total int

	var offset uint64 = 0
	if page > 0 {
		offset = (page - 1) * limit
	}

	err := q.DB.FluentSelect(func(query fluentsql.SelectBuilder) fluentsql.SelectBuilder {
		return query.
			From("users").
			Where(fluentsql.Eq{"deleted_at": nil}).
			Offset(offset).
			Limit(limit)
	}, &users, &total)

	// Return empty object and error.
	if err != nil {
		return users, total, err
	}

	// Return query result.
	return users, total, nil
}
