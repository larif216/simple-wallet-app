package config

import (
	"net/http"
	"simple-wallet-app/module/wallet/internal/handler"
)

func RegisterWalletHandlers(mux *http.ServeMux, config *WalletConfig) {
	h := handler.NewWalletHandler()
	mux.HandleFunc("/wallet/disburse", h.Disburse)
}
