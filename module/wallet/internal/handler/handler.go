package handler

import (
	"fmt"
	"net/http"
)

type WalletHandler struct {
}

func NewWalletHandler() *WalletHandler {
	return &WalletHandler{}
}

func (h *WalletHandler) Disburse(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
}
