package usecase

import (
	"context"
	"time"
	"user-service/internal/util"
	"user-service/internal/util/liberr"
	"user-service/internal/util/libvalidate"
	"user-service/module/auth/entity"

	"github.com/golang-jwt/jwt/v5"
)

type AuthUsecaseRepos struct {
	UserRepo UserRepository
}

type AuthUsecaseConfig struct {
	JWTSecret         string
	JWTHourExpiration int
}

type AuthUsecase struct {
	repos  *AuthUsecaseRepos
	config *AuthUsecaseConfig
}

func NewAuthUsecase(repos *AuthUsecaseRepos, config *AuthUsecaseConfig) *AuthUsecase {
	return &AuthUsecase{
		repos:  repos,
		config: config,
	}
}

func (a *AuthUsecase) Authentication(ctx context.Context, params *entity.AuthenticationUserRequest) (*entity.AuthenticationUser, error) {
	if err := libvalidate.Validator().Struct(params); err != nil {
		return nil, libvalidate.ResolveError(err, entity.ErrorCodeAuthInvalid)
	}

	user, err := a.repos.UserRepo.GetByUsername(ctx, params.Username)
	if err != nil {
		return nil, liberr.ResolveError(entity.ErrorAuthInvalid)
	}

	if user.Password != util.HashPlainTextWithSHA256(params.Password) {
		return nil, liberr.ResolveError(entity.ErrorAuthInvalid)
	}

	jwtSecret := []byte(a.config.JWTSecret)
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * time.Duration(a.config.JWTHourExpiration)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return nil, liberr.ResolveError(err)
	}

	return &entity.AuthenticationUser{
		AccessToken: tokenString,
	}, nil
}
