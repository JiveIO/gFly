package repository

// ====================================================================
// ======================== Repository factory ========================
// ====================================================================

// Repositories struct for collect all app repositories.
type Repositories struct {
	IRoleRepository
}

// Pool a repository pool to store all
var Pool = &Repositories{
	&RoleRepository{},
}
