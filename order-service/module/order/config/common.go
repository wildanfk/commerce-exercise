package config

import (
	"order-service/internal/util"
	"order-service/module/order/internal/repository"
	"order-service/module/order/internal/usecase"
	"time"

	rh "github.com/hashicorp/go-retryablehttp"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type OrderConfig struct {
	DB     *sqlx.DB    `ignored:"true"`
	Logger *zap.Logger `ignored:"true"`

	ProductServiceHost              string `envconfig:"PRODUCT_SERVICE_HOST" required:"true"`
	ProductServiceBasicAuthUsername string `envconfig:"PRODUCT_SERVICE_BASIC_AUTH_USERNAME" required:"true"`
	ProductServiceBasicAuthPassword string `envconfig:"PRODUCT_SERVICE_BASIC_AUTH_PASSWORD" required:"true"`

	WarehouseServiceHost              string `envconfig:"WAREHOUSE_SERVICE_HOST" required:"true"`
	WarehouseServiceBasicAuthUsername string `envconfig:"WAREHOUSE_SERVICE_BASIC_AUTH_USERNAME" required:"true"`
	WarehouseServiceBasicAuthPassword string `envconfig:"WAREHOUSE_SERVICE_BASIC_AUTH_PASSWORD" required:"true"`

	AuthServiceJWTSecret string `envconfig:"AUTH_SERVICE_JWT_SECRET" required:"true"`

	OrderExpirationTimeSecond int `envconfig:"SERVICE_ORDER_EXPIRATION_TIME_SECOND" required:"true"`
}

type repositorySet struct {
	orderRepository       *repository.OrderRepository
	orderDetailRepository *repository.OrderDetailRepository
	productRepository     *repository.ProductRepository
	warehouseRepository   *repository.WarehouseRepository
}

type usecaseSet struct {
	orderUsecase *usecase.OrderUsecase
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

func newRepositories(cfg *OrderConfig) (*repositorySet, error) {
	return &repositorySet{
		orderRepository:       repository.NewOrderRepository(cfg.DB),
		orderDetailRepository: repository.NewOrderDetailRepository(cfg.DB),
		warehouseRepository: repository.NewWarehouseRepository(
			repository.WarehouseConfiguration{
				ApiHost:           cfg.WarehouseServiceHost,
				BasicAuthUsername: cfg.WarehouseServiceBasicAuthUsername,
				BasicAuthPassword: cfg.WarehouseServiceBasicAuthPassword,
			},
			newServiceClient(),
		),
		productRepository: repository.NewProductRepository(
			repository.ProductConfiguration{
				ApiHost:           cfg.ProductServiceHost,
				BasicAuthUsername: cfg.ProductServiceBasicAuthUsername,
				BasicAuthPassword: cfg.ProductServiceBasicAuthPassword,
			},
			newServiceClient(),
		),
	}, nil
}

func newUsecase(cfg *OrderConfig, repositories *repositorySet) (*usecaseSet, error) {
	databaseTransactionHandler := util.NewDatabaseTransactionHandler(cfg.DB)

	return &usecaseSet{
		orderUsecase: usecase.NewOrderUsecase(&usecase.OrderUsecaseRepos{
			DatabaseTransactionHandler: databaseTransactionHandler,
			OrderRepo:                  repositories.orderRepository,
			OrderDetailRepo:            repositories.orderDetailRepository,
			WarehouseRepo:              repositories.warehouseRepository,
			ProductRepo:                repositories.productRepository,
		}, &usecase.OrderUsecaseConfig{
			OrderExpirationTimeSecond: cfg.OrderExpirationTimeSecond,
		}, cfg.Logger),
	}, nil
}
