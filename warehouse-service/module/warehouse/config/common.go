package config

import (
	"warehouse-service/internal/util"
	"warehouse-service/module/warehouse/internal/repository"
	"warehouse-service/module/warehouse/internal/usecase"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type WarehouseConfig struct {
	DB     *sqlx.DB    `ignored:"true"`
	Logger *zap.Logger `ignored:"true"`

	BasicAuthUsername string `envconfig:"SERVICE_BASIC_AUTH_USERNAME" required:"true"`
	BasicAuthPassword string `envconfig:"SERVICE_BASIC_AUTH_PASSWORD" required:"true"`
}

type repositorySet struct {
	warehouseRepository      *repository.WarehouseRepository
	warehouseStockRepository *repository.WarehouseStockRepository
}

type usecaseSet struct {
	warehouseUsecase      *usecase.WarehouseUsecase
	warehouseStockUsecase *usecase.WarehouseStockUsecase
}

func newRepositories(cfg *WarehouseConfig) (*repositorySet, error) {
	return &repositorySet{
		warehouseRepository:      repository.NewWarehouseRepository(cfg.DB),
		warehouseStockRepository: repository.NewWarehouseStockRepository(cfg.DB),
	}, nil
}

func newUsecase(cfg *WarehouseConfig, repositories *repositorySet) (*usecaseSet, error) {
	databaseTransactionHandler := util.NewDatabaseTransactionHandler(cfg.DB)

	return &usecaseSet{
		warehouseUsecase: usecase.NewWarehouseUsecase(&usecase.WarehouseUsecaseRepos{
			WarehouseRepo: repositories.warehouseRepository,
		}),
		warehouseStockUsecase: usecase.NewWarehouseStockUsecase(&usecase.WarehouseStockUsecaseRepos{
			DatabaseTransactionHandler: databaseTransactionHandler,
			WarehouseRepo:              repositories.warehouseRepository,
			WarehouseStockRepo:         repositories.warehouseStockRepository,
		}),
	}, nil
}
