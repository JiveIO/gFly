package repository

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
