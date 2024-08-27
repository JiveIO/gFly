package models

import (
	"database/sql"
	mb "github.com/gflydev/db"
	"time"
)

// ====================================================================
// ============================ Data Types ============================
// ====================================================================

// Define User Role Type enum

type RoleType uint

// Role name and id match with database.
const (
	RoleNA RoleType = iota
	RoleAdmin
	RoleModerator
	RoleMember
	RoleUser
	RoleGuest
)

var userRoleName = []string{
	"N/A",
	"admin",
	"moderator",
	"member",
	"user",
	"guest",
}

func (e RoleType) Name() string {
	return userRoleName[e]
}

func (e RoleType) Ordinal() int {
	return int(e)
}

func (e RoleType) Values() []string {
	return userRoleName
}

func (e RoleType) ByName(name string) RoleType {
	for k, v := range userRoleName {
		if name == v {
			return RoleType(k)
		}
	}

	return RoleNA
}

// ====================================================================
// ============================== Table ===============================
// ====================================================================

// TableRole Table name
const TableRole = "roles"

// Role struct to describe a role object.
type Role struct {
	// Table meta data
	MetaData mb.MetaData `db:"-" model:"table:roles"`

	// Table fields
	ID        int          `db:"id" model:"name:id; type:serial,primary"`
	Name      string       `db:"name" model:"name:name"`
	Slug      string       `db:"slug" model:"name:slug"`
	CreatedAt time.Time    `db:"created_at" model:"name:created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" model:"name:updated_at"`
}
