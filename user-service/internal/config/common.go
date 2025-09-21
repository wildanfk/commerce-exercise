package config

import (
	"fmt"
	"os"
	"time"

	authConfig "user-service/module/auth/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
	"go.uber.org/zap"
)

type ServiceConfig struct {
	GatewayHost       string         `envconfig:"GATEWAY_HOST" required:"true"`
	DatabaseConfig    DatabaseConfig `envconfig:"DB" required:"true"`
	JWTSecret         string         `envconfig:"JWT_SECRET" required:"true"`
	JWTHourExpiration int            `envconfig:"JWT_HOUR_EXPIRATION" required:"true"`

	DB     *sqlx.DB    `ignored:"true"`
	Logger *zap.Logger `ignored:"true"`
}

type DatabaseConfig struct {
	Driver          string `envconfig:"DRIVER" default:"mysql"`
	Host            string `envconfig:"HOST" default:"127.0.0.1"`
	Port            int    `envconfig:"PORT" default:"3306"`
	Username        string `envconfig:"USERNAME" required:"true"`
	Password        string `envconfig:"PASSWORD" required:"true"`
	Database        string `envconfig:"DATABASE" required:"true"`
	MaxLifetime     int    `envconfig:"MAXLIFETIME" default:"5"`
	MaxIdleLifetime int    `envconfig:"MAXIDLELIFETIME" default:"1"`
	MaxIdleConns    int    `envconfig:"MAXIDLECONNS" default:"25"`
	MaxOpenConns    int    `envconfig:"MAXOPENCONNS" default:"85"`
	QueryString     string `envconfig:"QUERYSTRING"`
}

func (c *DatabaseConfig) rWDataSourceName() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.QueryString,
	)
}

func newDatabase(cfg *ServiceConfig) *sqlx.DB {
	sqlxDB, err := sqlx.Connect(cfg.DatabaseConfig.Driver, cfg.DatabaseConfig.rWDataSourceName())
	if err != nil {
		panic(err)
	}
	sqlxDB.SetConnMaxLifetime(time.Minute * time.Duration(cfg.DatabaseConfig.MaxLifetime))
	sqlxDB.SetMaxIdleConns(cfg.DatabaseConfig.MaxIdleConns)
	sqlxDB.SetMaxOpenConns(cfg.DatabaseConfig.MaxOpenConns)

	return sqlxDB
}

func newLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger
}

func loadConfig() (*ServiceConfig, error) {
	cfg := new(ServiceConfig)

	// load from .env if exists
	if _, err := os.Stat(".env"); err == nil {
		if err := gotenv.Load(); err != nil {
			return nil, err
		}
	}

	// parse environment variable to config struct using prefix "service"
	if err := envconfig.Process("service", cfg); err != nil {
		return nil, err
	}

	cfg.DB = newDatabase(cfg)

	cfg.Logger = newLogger()

	return cfg, nil
}

func loadAuthConfig(serviceConfig *ServiceConfig) (*authConfig.AuthConfig, error) {
	cfg := new(authConfig.AuthConfig)

	if err := envconfig.Process("auth", cfg); err != nil {
		return nil, err
	}

	cfg.DB = serviceConfig.DB
	cfg.Logger = serviceConfig.Logger
	cfg.JWTSecret = serviceConfig.JWTSecret
	cfg.JWTHourExpiration = serviceConfig.JWTHourExpiration

	return cfg, nil
}
