package service

import (
	"auth-service/app/internal/model"
	"context"
)

type Auth interface {
	RegisterUser(
		ctx context.Context,
		userRequest *model.UserRegisterRequest,
	) (userResponse *model.UserTokenResponse, err error)
	LoginUser(
		ctx context.Context,
		userRequest *model.UserLoginRequest,
	) (userResponse *model.UserTokenResponse, err error)
	RefreshTokens(
		ctx context.Context,
		refreshToken string,
	) (userResponse *model.UserTokenResponse, err error)
}
