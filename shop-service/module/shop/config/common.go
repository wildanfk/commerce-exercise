package config

import (
	"shop-service/module/shop/internal/repository"
	"shop-service/module/shop/internal/usecase"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ShopConfig struct {
	DB     *sqlx.DB    `ignored:"true"`
	Logger *zap.Logger `ignored:"true"`

	BasicAuthUsername string `envconfig:"SERVICE_BASIC_AUTH_USERNAME" required:"true"`
	BasicAuthPassword string `envconfig:"SERVICE_BASIC_AUTH_PASSWORD" required:"true"`
}

type repositorySet struct {
	shopRepository *repository.ShopRepository
}

type usecaseSet struct {
	shopUsecase *usecase.ShopUsecase
}

func newRepositories(cfg *ShopConfig) (*repositorySet, error) {
	return &repositorySet{
		shopRepository: repository.NewShopRepository(cfg.DB),
	}, nil
}

func newUsecase(repositories *repositorySet) (*usecaseSet, error) {
	return &usecaseSet{
		shopUsecase: usecase.NewShopUsecase(&usecase.ShopUsecaseRepos{
			ShopRepo: repositories.shopRepository,
		}),
	}, nil
}
