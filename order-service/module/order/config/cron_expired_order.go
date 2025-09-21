package config

import (
	"order-service/internal/util/libcron"
	"order-service/module/order/internal/cron"
)

func NewCronExpiredOrder(cfg *OrderConfig) (*libcron.Cron, error) {
	repositories, err := newRepositories(cfg)
	if err != nil {
		return nil, err
	}

	usecases, err := newUsecase(cfg, repositories)
	if err != nil {
		return nil, err
	}

	cronHandler := cron.NewExpiredOrderCron(usecases.orderUsecase)

	return libcron.NewCron(libcron.Config{
		Name:        "CronOrderExpired",
		CronHandler: cronHandler,
		Logger:      cfg.Logger,
	}), nil
}
