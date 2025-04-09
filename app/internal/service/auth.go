package service

import (
	"auth-service/app/internal/exception"
	"auth-service/app/internal/lib"
	"auth-service/app/internal/model"
	"auth-service/app/internal/passhash"
	"auth-service/app/internal/repository"
	"auth-service/app/internal/token"
	"context"
	"log/slog"
	"time"
)

type AuthService struct {
	uow    repository.UnitOfWork
	logger *slog.Logger
}

func NewAuthService(uow repository.UnitOfWork, logger *slog.Logger) *AuthService {
	return &AuthService{uow: uow, logger: logger}
}

func (s *AuthService) RegisterUser(
	ctx context.Context,
	userRequest *model.UserRegisterRequest,
) (userResponse *model.UserTokenResponse, err error) {
	err = s.uow.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			s.uow.Rollback()
		} else {
			err = s.uow.Commit()
			if err != nil {
				userResponse = nil
			}
		}
	}()

	// create user
	userModel, err := userRequest.ToUser()
	if err != nil {
		s.logger.Error("Model conversion error", "error", err)
		return nil, exception.ErrInternal
	}

	userRepo := s.uow.UserRepository()
	err = userRepo.CreateOne(ctx, userModel)
	if err != nil {
		return nil, err
	}

	userRoleModel := &model.UserRole{UserID: userModel.ID, RoleName: lib.User}
	userRoleRepo := s.uow.UserRoleRepository()
	err = userRoleRepo.CreateOne(ctx, userRoleModel)
	if err != nil {
		return nil, err
	}

	// create tokens
	refreshToken := token.CreateRefresh()
	refresSessionModel := &model.RefreshSession{
		RefreshToken: refreshToken,
		UserID:       userModel.ID,
	}

	refreshSessionRepo := s.uow.RefreshSessionRepository()
	err = refreshSessionRepo.CreateOne(ctx, refresSessionModel)
	if err != nil {
		return nil, err
	}

	access_token, err := token.CreateAccess(userModel.ID, lib.User)
	if err != nil {
		s.logger.Error("Access token create failed")
		return nil, err
	}

	userResponse = &model.UserTokenResponse{
		RefreshToken: refreshToken,
		AccessToken:  access_token,
	}
	return userResponse, nil
}

func (s *AuthService) LoginUser(
	ctx context.Context,
	userRequest *model.UserLoginRequest,
) (userResponse *model.UserTokenResponse, err error) {
	err = s.uow.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			s.uow.Rollback()
		} else {
			err = s.uow.Commit()
			if err != nil {
				userResponse = nil
			}
		}
	}()

	// user exist check
	var userModel *model.User
	userRepo := s.uow.UserRepository()
	if isEmail := lib.ValidateEmail(userRequest.Login); isEmail {
		userModel, err = userRepo.GetByEmail(ctx, userRequest.Login)
		if err != nil {
			return nil, err
		}
	} else {
		userModel, err = userRepo.GetByUsername(ctx, userRequest.Login)
		if err != nil {
			return nil, err
		}
	}

	if !(passhash.Verify(userModel.Password, userRequest.Password)) {
		return nil, exception.ErrVerify
	}

	// create tokens
	refreshToken := token.CreateRefresh()
	sessionModel := &model.RefreshSession{RefreshToken: refreshToken, UserID: userModel.ID}

	sessionRepo := s.uow.RefreshSessionRepository()
	err = sessionRepo.CreateOne(ctx, sessionModel)
	if err != nil {
		return nil, err
	}

	userRoleRepo := s.uow.UserRoleRepository()
	userRole, err := userRoleRepo.GetRoleByUserID(ctx, userModel.ID)
	if err != nil {
		return nil, err
	}

	accessToken, err := token.CreateAccess(userModel.ID, userRole)
	if err != nil {
		return nil, err
	}

	userResponse = &model.UserTokenResponse{AccessToken: accessToken, RefreshToken: refreshToken}
	return userResponse, nil
}

func (s *AuthService) RefreshTokens(
	ctx context.Context,
	refreshToken string,
) (userResponse *model.UserTokenResponse, err error) {
	err = s.uow.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			s.uow.Rollback()
		} else {
			err = s.uow.Commit()
			if err != nil {
				userResponse = nil
			}
		}
	}()

	sessionRepo := s.uow.RefreshSessionRepository()
	sessionModel, err := sessionRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if sessionModel.ExpireAt.After(time.Now().UTC()) {
		return nil, exception.ErrTokenExpired
	}

	refreshTokenNew := token.CreateRefresh()
	err = sessionRepo.UpdateTokenByID(ctx, sessionModel.ID, refreshTokenNew)
	if err != nil {
		return nil, err
	}

	userRoleRepo := s.uow.UserRoleRepository()
	userRole, err := userRoleRepo.GetRoleByUserID(ctx, sessionModel.UserID)
	if err != nil {
		return nil, err
	}

	accessToken, err := token.CreateAccess(sessionModel.UserID, userRole)
	if err != nil {
		return nil, err
	}

	userResponse = &model.UserTokenResponse{
		RefreshToken: refreshTokenNew,
		AccessToken:  accessToken,
	}
	return userResponse, nil
}
