package repository

import (
	"auth-service/app/internal/lib"
	"auth-service/app/internal/model"
)

type UserRepository interface {
	CreateOne(user *model.User) error
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetCountByID(id uint) (int, error)
}

type RefreshSessionRepository interface {
	CreateOne(refreshSession *model.RefreshSession) error
	GetByToken(refreshToken string) (*model.RefreshSession, error)
	UpdateTokenByID(id uint, refreshToken string) error
}

type UserRoleRepository interface {
	CreateOne(userRole *model.UserRole) error
	GetRoleByUserID(userID uint) (lib.Role, error)
}

type UnitOfWork interface {
	Begin() error
	Commit() error
	Rollback() error
	UserRepository() UserRepository
	RefreshSessionRepository() RefreshSessionRepository
	UserRoleRepository() UserRoleRepository
}
