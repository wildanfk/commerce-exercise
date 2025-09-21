package config

import (
	"user-service/module/auth/internal/repository"
	"user-service/module/auth/internal/usecase"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type AuthConfig struct {
	DB                *sqlx.DB    `ignored:"true"`
	JWTSecret         string      `ignored:"true"`
	JWTHourExpiration int         `ignored:"true"`
	Logger            *zap.Logger `ignored:"true"`
}

type repositorySet struct {
	userRepository *repository.UserRepository
}

type usecaseSet struct {
	authUsecase *usecase.AuthUsecase
}

func newRepositories(cfg *AuthConfig) (*repositorySet, error) {
	return &repositorySet{
		userRepository: repository.NewUserRepository(cfg.DB),
	}, nil
}

func newUsecase(cfg *AuthConfig, repositories *repositorySet) (*usecaseSet, error) {
	return &usecaseSet{
		authUsecase: usecase.NewAuthUsecase(&usecase.AuthUsecaseRepos{
			UserRepo: repositories.userRepository,
		},
			&usecase.AuthUsecaseConfig{
				JWTSecret:         cfg.JWTSecret,
				JWTHourExpiration: cfg.JWTHourExpiration,
			}),
	}, nil
}
