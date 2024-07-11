package repository

import (
	"app/app/domain/models"
	mb "app/core/fluentmodel" // Model builder
	qb "app/core/fluentsql"   // Query builder
	"github.com/google/uuid"
)

// UserRepository struct for queries from a User model.
// The struct is an implementation of interface IUserRepository
type UserRepository struct {
}

// GetUserByID query for getting one User by given ID.
func (q *UserRepository) GetUserByID(id uuid.UUID) (models.User, error) {
	// DB Model instance
	db := mb.Instance()
	// Error variable
	var err error

	// Define User variable.
	user := models.User{}
	err = db.Where("id", qb.Eq, id).
		First(&user)

	// Return query result.
	return user, err
}

// GetUserByEmail query for getting one User by given Email.
func (q *UserRepository) GetUserByEmail(email string) (models.User, error) {
	// DB Model instance
	db := mb.Instance()
	// Error variable
	var err error

	// Define User variable.
	user := models.User{}
	err = db.Where("email", qb.Eq, email).
		First(&user)

	// Return query result.
	return user, err
}

// CreateUser a query for creating a new user by given user data.
func (q *UserRepository) CreateUser(u *models.User) error {
	// DB Model instance
	db := mb.Instance()

	// Defer a rollback in case anything fails.
	db.Begin()
	defer func(db *mb.DBModel) {
		_ = db.Rollback()
	}(db)

	// Error variable
	var err error

	// Create new user
	err = db.Create(u)
	if err != nil {
		return err
	}

	// DB commit
	err = db.Commit()

	// This query returns nothing.
	return err
}

// UpdateUser a query for updating a user by given user data.
func (q *UserRepository) UpdateUser(u *models.User) error {
	// DB Model instance
	db := mb.Instance()

	// Defer a rollback in case anything fails.
	db.Begin()
	defer func(db *mb.DBModel) {
		_ = db.Rollback()
	}(db)

	// Error variable
	var err error

	// Create new user
	err = db.Update(u)
	if err != nil {
		return err
	}

	// DB commit
	err = db.Commit()

	// This query returns nothing.
	return err
}

// DeleteUser a query for updating a user by given user data.
func (q *UserRepository) DeleteUser(u *models.User) error {
	// DB Model instance
	db := mb.Instance()

	// Defer a rollback in case anything fails.
	db.Begin()
	defer func(db *mb.DBModel) {
		_ = db.Rollback()
	}(db)

	// Error variable
	var err error

	// Create new user
	err = db.Delete(u)
	if err != nil {
		return err
	}

	// DB commit
	err = db.Commit()

	// This query returns nothing.
	return err
}

// SelectUser a query for updating a user by given user data.
func (q *UserRepository) SelectUser(page, limit int) ([]*models.User, int, error) {
	// DB Model instance
	db := mb.Instance()
	// Error variable
	var err error

	// Define User variable.
	var users []*models.User
	var total int

	var offset = 0
	if page > 0 {
		offset = (page - 1) * limit
	}

	total, err = db.Limit(limit, offset).
		Find(users)

	// Return query result.
	return users, total, err
}
