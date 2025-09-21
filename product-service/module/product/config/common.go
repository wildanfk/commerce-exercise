package config

import (
	"product-service/module/product/internal/repository"
	"product-service/module/product/internal/usecase"
	"time"

	rh "github.com/hashicorp/go-retryablehttp"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type ProductConfig struct {
	DB     *sqlx.DB    `ignored:"true"`
	Logger *zap.Logger `ignored:"true"`

	BasicAuthUsername string `envconfig:"SERVICE_BASIC_AUTH_USERNAME" required:"true"`
	BasicAuthPassword string `envconfig:"SERVICE_BASIC_AUTH_PASSWORD" required:"true"`

	WarehouseServiceHost              string `envconfig:"WAREHOUSE_SERVICE_HOST" required:"true"`
	WarehouseServiceBasicAuthUsername string `envconfig:"WAREHOUSE_SERVICE_BASIC_AUTH_USERNAME" required:"true"`
	WarehouseServiceBasicAuthPassword string `envconfig:"WAREHOUSE_SERVICE_BASIC_AUTH_PASSWORD" required:"true"`

	ShopServiceHost              string `envconfig:"SHOP_SERVICE_HOST" required:"true"`
	ShopServiceBasicAuthUsername string `envconfig:"SHOP_SERVICE_BASIC_AUTH_USERNAME" required:"true"`
	ShopServiceBasicAuthPassword string `envconfig:"SHOP_SERVICE_BASIC_AUTH_PASSWORD" required:"true"`
}

type repositorySet struct {
	productRepository   *repository.ProductRepository
	warehouseRepository *repository.WarehouseRepository
	shopRepository      *repository.ShopRepository
}

type usecaseSet struct {
	productUsecase *usecase.ProductUsecase
}

func newServiceClient() *rh.Client {
	client := rh.NewClient()
	client.Logger = nil // disable internal logging
	client.HTTPClient.Timeout = time.Duration(10) * time.Second
	client.RetryMax = 3
	client.Backoff = rh.DefaultBackoff
	client.CheckRetry = rh.DefaultRetryPolicy

	return client
}

func newRepositories(cfg *ProductConfig) (*repositorySet, error) {
	return &repositorySet{
		productRepository: repository.NewProductRepository(cfg.DB),
		warehouseRepository: repository.NewWarehouseRepository(
			repository.WarehouseConfiguration{
				ApiHost:           cfg.WarehouseServiceHost,
				BasicAuthUsername: cfg.WarehouseServiceBasicAuthUsername,
				BasicAuthPassword: cfg.WarehouseServiceBasicAuthPassword,
			},
			newServiceClient(),
		),
		shopRepository: repository.NewShopRepository(
			repository.ShopConfiguration{
				ApiHost:           cfg.ShopServiceHost,
				BasicAuthUsername: cfg.ShopServiceBasicAuthUsername,
				BasicAuthPassword: cfg.ShopServiceBasicAuthPassword,
			},
			newServiceClient(),
		),
	}, nil
}

func newUsecase(repositories *repositorySet) (*usecaseSet, error) {
	return &usecaseSet{
		productUsecase: usecase.NewProductUsecase(&usecase.ProductUsecaseRepos{
			ProductRepo:   repositories.productRepository,
			WarehouseRepo: repositories.warehouseRepository,
			ShopRepo:      repositories.shopRepository,
		}),
	}, nil
}
