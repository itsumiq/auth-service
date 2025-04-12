package postgres

import (
	"auth-service/app/internal/exception"
	"auth-service/app/internal/lib"
	"auth-service/app/internal/model"
	"database/sql"
	"errors"
	"log/slog"

	"golang.org/x/net/context"
)

type userRoleRepository struct {
	conn   Conn
	logger *slog.Logger
}

func NewUserRoleRepository(db Conn, logger *slog.Logger) *userRoleRepository {
	return &userRoleRepository{conn: db, logger: logger}
}

func (r *userRoleRepository) CreateOne(ctx context.Context, userRole *model.UserRole) error {
	query := `
	INSERT INTO users_roles (user_id, role_name)
	VALUES ($1, $2)
	`
	_, err := r.conn.ExecContext(ctx, query, userRole.UserID, userRole.RoleName)
	if err != nil {
		r.logger.Error("Database insert failed", "error", err)
		return err
	}

	return nil
}

func (r *userRoleRepository) GetRoleByUserID(ctx context.Context, userID uint) (lib.Role, error) {
	var role lib.Role
	query := `
	SELECT role_name FROM users_roles
	WHERE user_id = $1
	`

	row := r.conn.QueryRowxContext(ctx, query, userID)
	if err := row.Scan(&role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", exception.ErrNotFound
		}
		r.logger.Error("Database select failed", "error", err)
		return "", exception.ErrInternal
	}

	return role, nil
}
