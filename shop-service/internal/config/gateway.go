package config

import (
	"net/http"
	shopConfig "shop-service/module/shop/config"

	"github.com/gorilla/mux"
)

func NewGatewayServer() (*http.Server, *ServiceConfig, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, nil, err
	}

	authCfg, err := loadAuthConfig(cfg)
	if err != nil {
		return nil, nil, err
	}

	serverMux := mux.NewRouter()
	shopConfig.RegisterGatewayHandler(serverMux, authCfg)

	return &http.Server{
		Addr:    cfg.GatewayHost,
		Handler: serverMux,
	}, cfg, nil
}
