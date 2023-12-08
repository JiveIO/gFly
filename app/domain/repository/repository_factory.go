package repository

import "app/core/db"

// Repositories struct for collect all app repositories.
type Repositories struct {
	*UserRepository
	*RoleRepository
}

// Pool a repository pool to store all
var Pool = &Repositories{
	&UserRepository{DB: db.Instance()},
	&RoleRepository{DB: db.Instance()},
}
