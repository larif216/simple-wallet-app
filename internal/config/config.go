package config

import (
	"fmt"
	"net/http"
	"os"

	wconfig "simple-wallet-app/module/wallet/config"

	"github.com/kelseyhightower/envconfig"
	"github.com/subosito/gotenv"
)

type ServiceConfig struct {
	Host           string         `envconfig:"HOST"`
	DatabaseConfig DatabaseConfig `envconfig:"DB"`
}

type DatabaseConfig struct {
	Driver      string `envconfig:"DRIVER"`
	Host        string `envconfig:"HOST"`
	Port        string `envconfig:"PORT"`
	Username    string `envconfig:"USERNAME"`
	Password    string `envconfig:"PASSWORD"`
	Database    string `envconfig:"DATABASE"`
	QueryString string `envconfig:"QUERY_STRING"`
}

type HttpServer struct {
	HTTPServer *http.Server
	Config     *ServiceConfig
}

func NewHttpServer() (*HttpServer, error) {
	svcCfg, err := loadServiceConfig()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	wcfg := wconfig.LoadWalletConfig()
	wcfg.Database = newDatabase(svcCfg.DatabaseConfig)
	wconfig.RegisterWalletHandlers(mux, wcfg)

	hs := &http.Server{
		Addr:    svcCfg.Host,
		Handler: mux,
	}

	return &HttpServer{
		HTTPServer: hs,
		Config:     &svcCfg,
	}, nil
}

func (c *DatabaseConfig) RWDataSourceName() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
		c.QueryString,
	)
}

func loadServiceConfig() (ServiceConfig, error) {
	var cfg ServiceConfig

	if _, err := os.Stat(".env"); err == nil {
		if err := gotenv.Load(); err != nil {
			return cfg, err
		}
	}

	err := envconfig.Process("wallet_service", &cfg)

	return cfg, err
}
