package usecase

import (
	"simple-wallet-app/module/wallet/entity"

	"github.com/shopspring/decimal"
)

func (uc *WalletUsecase) Disburse(req entity.DisburseRequest) (*entity.DisburseResponse, error) {
	tx, err := uc.dbTransaction.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	user, err := uc.userRepo.GetByID(tx, req.UserID)
	if err != nil {
		return nil, err
	}

	availableBalance := user.Balance.Sub(user.PendingBalance)
	if availableBalance.LessThan(req.Amount) {
		return nil, entity.ErrInsufficientBalance
	}

	disbursement := &entity.Disbursement{
		UserID: req.UserID,
		Amount: req.Amount,
	}

	err = uc.disbursementRepo.Create(tx, disbursement)
	if err != nil {
		return nil, err
	}

	err = uc.userRepo.HoldBalance(tx, req.UserID, req.Amount)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	status, err := uc.process(disbursement.UserID, disbursement.ID, disbursement.Amount)
	if err != nil {
		return &entity.DisburseResponse{
			DisbursementID:         disbursement.ID,
			DisbursementStatusEnum: status,
			Message:                err.Error(),
		}, err
	}

	var message string
	switch status {
	case entity.DisbursementStatusSuccess:
		message = "Disbursement successfully processed"
	case entity.DisbursementStatusFailed:
		message = "Disbursement failed due to processing error"
	default:
		message = "Disbursement processing is pending"
	}

	return &entity.DisburseResponse{
		DisbursementID:         disbursement.ID,
		DisbursementStatusEnum: status,
		Message:                message,
	}, nil
}

func (uc *WalletUsecase) process(userID, disbursementID uint32, amount decimal.Decimal) (entity.DisbursementStatusEnum, error) {
	status := entity.DisbursementStatusPending

	tx, err := uc.dbTransaction.Begin()
	if err != nil {
		return status, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// mocking purposes, to simulate success, pending, and failed scenario
	lastDigit := amount.IntPart() % 10

	switch lastDigit {
	case 1:
		// failed case
		err = uc.userRepo.ReleaseBalance(tx, userID, amount, true)
		if err != nil {
			return status, err
		}

		status = entity.DisbursementStatusFailed

		err = uc.disbursementRepo.UpdateStatus(tx, disbursementID, status)
		if err != nil {
			return status, err
		}
	case 2:
		// pending case
		status = entity.DisbursementStatusPending
	default:
		// success case
		err = uc.userRepo.ReleaseBalance(tx, userID, amount, false)
		if err != nil {
			return status, err
		}

		status = entity.DisbursementStatusSuccess

		err = uc.disbursementRepo.UpdateStatus(tx, disbursementID, status)
		if err != nil {
			return status, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return status, err
	}

	return status, nil
}
