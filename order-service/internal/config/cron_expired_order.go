package config

import (
	"order-service/internal/util/libcron"
	orderConfig "order-service/module/order/config"
)

func NewCronExpiredOrder() (*libcron.Cron, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}

	contentCfg, err := loadOrderConfig(cfg)
	if err != nil {
		return nil, err
	}

	return orderConfig.NewCronExpiredOrder(contentCfg)
}
