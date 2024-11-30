package config

import (
	"simple-wallet-app/internal/sqlutil"
	"simple-wallet-app/module/wallet/internal/repository"
	"simple-wallet-app/module/wallet/internal/usecase"
)

func NewWalletUsecase(cfg *WalletConfig) *usecase.WalletUsecase {
	var userRepo usecase.UserRepository
	var disbursementRepo usecase.DisbursementRespository

	userRepo = repository.NewUserRepository(cfg.Database)
	disbursementRepo = repository.NewDisbursementRepository(cfg.Database)

	dbTransaction := sqlutil.NewDatabaseTransactionCreator(cfg.Database)

	uc := usecase.NewWalletUseCase(
		userRepo,
		disbursementRepo,
		dbTransaction,
	)

	return uc
}
