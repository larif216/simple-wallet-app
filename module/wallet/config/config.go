package config

import "database/sql"

type WalletConfig struct {
	Database *sql.DB
}

func LoadWalletConfig() *WalletConfig {
	cfg := new(WalletConfig)
	return cfg
}
