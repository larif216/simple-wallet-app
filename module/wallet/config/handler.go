package config

import (
	"net/http"
	"simple-wallet-app/module/wallet/internal/handler"
)

func RegisterWalletHandlers(mux *http.ServeMux, config *WalletConfig) {
	uc := NewWalletUsecase(config)
	h := handler.NewWalletHandler(uc)
	mux.HandleFunc("/wallet/disburse", h.Disburse)
}
