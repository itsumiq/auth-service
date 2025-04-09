package model

import role "auth-service/app/internal/lib"

type UserRole struct {
	UserID   uint      `db:"user_id"`
	RoleName role.Role `db:"role_name"`
}
