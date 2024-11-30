package usecase

import (
	"simple-wallet-app/internal/sqlutil"
	"simple-wallet-app/module/wallet/entity"

	"github.com/shopspring/decimal"
)

type UserRepository interface {
	GetByID(tx sqlutil.DatabaseTransaction, userID uint32) (*entity.User, error)
	HoldBalance(tx sqlutil.DatabaseTransaction, userID uint32, amount decimal.Decimal) error
	ReleaseBalance(tx sqlutil.DatabaseTransaction, userID uint32, amount decimal.Decimal, withReversal bool) error
}

type DisbursementRespository interface {
	Create(tx sqlutil.DatabaseTransaction, params *entity.Disbursement) error
	UpdateStatus(tx sqlutil.DatabaseTransaction, disbursementID uint32, status entity.DisbursementStatusEnum) error
}

type WalletUsecases interface {
	Disburse(req entity.DisburseRequest) (*entity.DisburseResponse, error)
}
