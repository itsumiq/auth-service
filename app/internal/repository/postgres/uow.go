package postgres

import (
	"auth-service/app/internal/exception"
	"auth-service/app/internal/repository"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type unitOfWork struct {
	db     *sqlx.DB
	tx     *sqlx.Tx
	logger *slog.Logger
}

func NewUnitOfWork(db *sqlx.DB, logger *slog.Logger) *unitOfWork {
	return &unitOfWork{db: db, logger: logger}
}

func (uow *unitOfWork) Begin() error {
	tx, err := uow.db.Beginx()
	if err != nil {
		uow.logger.Error("Database start transaction failed", "error", err)
		return exception.ErrInternal
	}
	uow.tx = tx

	return nil
}

func (uow *unitOfWork) Commit() error {
	if err := uow.tx.Commit(); err != nil {
		uow.logger.Error("Database commit transaction failed", "error", err)
		return exception.ErrInternal
	}

	return nil
}

func (uow *unitOfWork) Rollback() error {
	if err := uow.tx.Rollback(); err != nil {
		uow.logger.Error("Database rollbackk transaction failed", "error", err)
		return exception.ErrInternal
	}

	return nil
}

func (uow *unitOfWork) UserRepository() repository.UserRepository {
	if uow.tx != nil {
		return NewUserRepository(uow.tx, uow.logger)
	}
	return NewUserRepository(uow.db, uow.logger)
}

func (uow *unitOfWork) RefreshSessionRepository() repository.RefreshSessionRepository {
	if uow.tx != nil {
		return NewRefreshSessionRepository(uow.tx, uow.logger)
	}
	return NewRefreshSessionRepository(uow.db, uow.logger)
}

func (uow *unitOfWork) UserRoleRepository() repository.UserRoleRepository {
	if uow.tx != nil {
		return NewUserRoleRepository(uow.tx, uow.logger)
	}
	return NewUserRoleRepository(uow.db, uow.logger)
}
