package postgres

import (
	"auth-service/app/internal/exception"
	"auth-service/app/internal/model"
	"context"
	"database/sql"
	"errors"
	"log/slog"
)

type refreshSessionRepository struct {
	conn   Conn
	logger *slog.Logger
}

func NewRefreshSessionRepository(conn Conn, logger *slog.Logger) *refreshSessionRepository {
	return &refreshSessionRepository{conn: conn, logger: logger}
}

func (r *refreshSessionRepository) CreateOne(ctx context.Context, refreshSession *model.RefreshSession) error {
	query := `
	INSERT INTO refresh_sessions (user_id, refresh_token)
	VALUES ($1, $2)
	`
	_, err := r.conn.ExecContext(
		ctx,
		query,
		refreshSession.UserID,
		refreshSession.RefreshToken,
	)
	if err != nil {
		r.logger.Error("Database insert failed", "error", err)
		return exception.ErrInternal
	}

	return nil
}

func (r *refreshSessionRepository) GetByToken(ctx context.Context, refreshToken string) (*model.RefreshSession, error) {
	refreshSession := &model.RefreshSession{}
	query := `
	SELECT * from refresh_sessions
	WHERE refresh_token = $1
	`

	row := r.conn.QueryRowxContext(ctx, query, refreshToken)
	if err := row.StructScan(refreshSession); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, exception.ErrNotFound
		}
		r.logger.Error("Database select failed", "error", err)
		return nil, exception.ErrInternal
	}

	return refreshSession, nil
}

func (r *refreshSessionRepository) UpdateTokenByID(ctx context.Context, id uint, refreshToken string) error {
	query := `
	UPDATE refresh_sessions
	SET refresh_token = $1
	WHERE id = $2
	`

	result, err := r.conn.ExecContext(ctx, query, refreshToken, id)
	if err != nil {
		r.logger.Error("Database update failed", "error", err)
		return exception.ErrInternal
	}
	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return exception.ErrNotFound
	}

	return nil
}
