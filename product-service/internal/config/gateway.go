package config

import (
	"net/http"
	productConfig "product-service/module/product/config"

	"github.com/gorilla/mux"
)

func NewGatewayServer() (*http.Server, *ServiceConfig, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, nil, err
	}

	authCfg, err := loadProductConfig(cfg)
	if err != nil {
		return nil, nil, err
	}

	serverMux := mux.NewRouter()
	productConfig.RegisterGatewayHandler(serverMux, authCfg)

	return &http.Server{
		Addr:    cfg.GatewayHost,
		Handler: serverMux,
	}, cfg, nil
}
