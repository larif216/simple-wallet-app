package usecase

import (
	"simple-wallet-app/internal/sqlutil"
)

type WalletUsecase struct {
	userRepo         UserRepository
	disbursementRepo DisbursementRespository

	dbTransaction sqlutil.DatabaseTransactionCreator
}

func NewWalletUsecase(
	userRepo UserRepository,
	disbursementRepo DisbursementRespository,
	dbTransaction sqlutil.DatabaseTransactionCreator,
) *WalletUsecase {
	return &WalletUsecase{
		userRepo:         userRepo,
		disbursementRepo: disbursementRepo,
		dbTransaction:    dbTransaction,
	}
}
