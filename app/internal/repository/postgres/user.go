package postgres

import (
	"auth-service/app/internal/exception"
	"auth-service/app/internal/model"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Conn interface {
	sqlx.ExtContext
	GetContext(ctx context.Context, dest any, query string, args ...any) error
}

type userRepository struct {
	conn   Conn
	logger *slog.Logger
	ctx    context.Context
}

func NewUserRepository(db Conn, logger *slog.Logger, ctx context.Context) *userRepository {
	return &userRepository{conn: db, logger: logger, ctx: ctx}
}

func (r *userRepository) CreateOne(user *model.User) error {
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err := r.conn.GetContext(r.ctx, &user.ID, query, user.Username, user.Email, user.Password)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return exception.ErrDuplicateEntry
		}
		r.logger.Error("Database insert failed", "error", err)
		return err
	}

	return nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user = &model.User{}
	query := `
	SELECT * FROM users
	WHERE username = $1
	`

	row := r.conn.QueryRowxContext(r.ctx, query, username)
	if err := row.StructScan(user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exception.ErrNotFound
		}
		r.logger.Error("Database select failed", "error", err)
		return nil, exception.ErrInternal
	}

	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user = &model.User{}
	query := `
	SELECT * FROM users
	WHERE email = $1
	`

	row := r.conn.QueryRowxContext(r.ctx, query, email)
	if err := row.StructScan(user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exception.ErrNotFound
		}
		r.logger.Error("Database select failed", "error", err)
		return nil, exception.ErrInternal
	}

	return user, nil
}

func (r *userRepository) GetCountByID(id uint) (int, error) {
	var count int
	query := `
	SELECT COUNT(*) FROM users
	WHERE id = $1
	`

	row := r.conn.QueryRowxContext(r.ctx, query, id)
	if err := row.StructScan(&count); err != nil {
		r.logger.Error("Database scan result failed", "error", err)
		return count, exception.ErrInternal
	}

	return count, nil
}
