package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple-wallet-app/internal/util"
	"simple-wallet-app/module/wallet/entity"
	"simple-wallet-app/module/wallet/internal/usecase"
)

type WalletHandler struct {
	walletUsecase usecase.WalletUsecases
}

func NewWalletHandler(uc usecase.WalletUsecases) *WalletHandler {
	return &WalletHandler{
		walletUsecase: uc,
	}
}

func (h *WalletHandler) Disburse(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		util.WriteHTTPResponse(w, map[string]string{
			"message": "Invalid request method",
		}, http.StatusMethodNotAllowed)
		return
	}

	var req entity.DisburseRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		util.WriteHTTPResponse(w, map[string]string{
			"message": fmt.Sprintf("Error parsing request body: %v", err),
		}, http.StatusBadRequest)
		return
	}

	resp, err := h.walletUsecase.Disburse(req)
	if err != nil {
		switch err {
		case entity.ErrUserNotFound:
			util.WriteHTTPResponse(w, map[string]string{
				"message": fmt.Sprintf("%v", err),
			}, http.StatusNotFound)
			return
		case entity.ErrDisbursementNotFound:
			util.WriteHTTPResponse(w, map[string]string{
				"message": fmt.Sprintf("%v", err),
			}, http.StatusNotFound)
			return
		case entity.ErrInsufficientBalance:
			util.WriteHTTPResponse(w, map[string]string{
				"message": fmt.Sprintf("%v", err),
			}, http.StatusBadRequest)
			return
		default:
			util.WriteHTTPResponse(w, map[string]string{
				"message": fmt.Sprintf("Error processing disbursement: %v", err),
			}, http.StatusInternalServerError)
			return
		}
	}

	response := map[string]interface{}{
		"disbursement_id":     resp.DisbursementID,
		"disbursement_status": resp.DisbursementStatusEnum.String(),
		"message":             resp.Message,
	}

	util.WriteHTTPResponse(w, response, http.StatusOK)
}
