package repository

import (
	"auth-service/app/internal/lib"
	"auth-service/app/internal/model"
	"context"
)

type UserRepository interface {
	CreateOne(ctx context.Context, user *model.User) error
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetCountByID(ctx context.Context, id uint) (int, error)
}

type RefreshSessionRepository interface {
	CreateOne(ctx context.Context, refreshSession *model.RefreshSession) error
	GetByToken(ctx context.Context, refreshToken string) (*model.RefreshSession, error)
	UpdateTokenByID(ctx context.Context, id uint, refreshToken string) error
}

type UserRoleRepository interface {
	CreateOne(ctx context.Context, userRole *model.UserRole) error
	GetRoleByUserID(ctx context.Context, userID uint) (lib.Role, error)
}

type UnitOfWork interface {
	Begin() error
	Commit() error
	Rollback() error
	UserRepository() UserRepository
	RefreshSessionRepository() RefreshSessionRepository
	UserRoleRepository() UserRoleRepository
}
