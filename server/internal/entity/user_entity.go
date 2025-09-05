package entity

import (
	"fmt"
	"time"
)

type UserRole string

const (
	UserRoleEmployee UserRole = "employee"
	UserRoleManager  UserRole = "manager"
)

type User struct {
	ID           uint64    `db:"id"`
	Email        string    `db:"email"`
	Name         string    `db:"name"`
	PasswordHash string    `db:"password_hash"`
	Role         UserRole  `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
}

func ParseUserRole(str string) (UserRole, error) {
	switch str {
	case "employee":
		return UserRoleEmployee, nil
	case "manager":
		return UserRoleManager, nil
	default:
		return "", fmt.Errorf("invalid user role = %s", str)
	}
}

type UserSimple struct {
	ID    uint64 `db:"id"`
	Email string `db:"email"`
	Name  string `db:"name"`
}
