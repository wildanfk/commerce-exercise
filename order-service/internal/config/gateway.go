package config

import (
	"net/http"
	orderConfig "order-service/module/order/config"

	"github.com/gorilla/mux"
)

func NewGatewayServer() (*http.Server, *ServiceConfig, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, nil, err
	}

	authCfg, err := loadOrderConfig(cfg)
	if err != nil {
		return nil, nil, err
	}

	serverMux := mux.NewRouter()
	orderConfig.RegisterGatewayHandler(serverMux, authCfg)

	return &http.Server{
		Addr:    cfg.GatewayHost,
		Handler: serverMux,
	}, cfg, nil
}
