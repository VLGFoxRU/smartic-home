package domain

import "time"

type UserRole string

const (
	RoleOwner      UserRole = "owner"
	RoleAdmin      UserRole = "admin"
	RoleController UserRole = "controller"
	RoleViewer     UserRole = "viewer"
)

type User struct {
	id           string
	username     string
	email        string
	passwordHash string
	role         UserRole
	createdAt    time.Time
}

func NewUser(id, username, email, passwordHash string, role UserRole) *User {
	return &User{
		id:           id,
		username:     username,
		email:        email,
		passwordHash: passwordHash,
		role:         role,
		createdAt:    time.Now(),
	}
}

func (u *User) ID() string           { return u.id }
func (u *User) Username() string     { return u.username }
func (u *User) Email() string        { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) Role() UserRole       { return u.role }
func (u *User) CreatedAt() time.Time { return u.createdAt }