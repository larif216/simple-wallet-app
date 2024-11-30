package repository

import (
	"database/sql"
	"simple-wallet-app/internal/sqlutil"
	"simple-wallet-app/module/wallet/entity"
	"time"
)

type DisbursementRepository struct {
	db *sql.DB
}

func NewDisbursementRepository(db *sql.DB) *DisbursementRepository {
	return &DisbursementRepository{
		db: db,
	}
}

func (r *DisbursementRepository) Create(tx sqlutil.DatabaseTransaction, params *entity.Disbursement) error {
	query := "INSERT INTO disbursements (user_id, amount, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	now := time.Now()

	var id uint32
	err := tx.QueryRow(
		query,
		params.UserID,
		params.Amount,
		entity.DisbursementStatusPending,
		now,
		now,
	).Scan(&id)
	if err != nil {
		return err
	}

	params.ID = id
	params.Status = entity.DisbursementStatusPending
	params.CreatedAt = now
	params.UpdatedAt = now

	return nil
}

func (r *DisbursementRepository) UpdateStatus(tx sqlutil.DatabaseTransaction, disbursementID uint32, status entity.DisbursementStatusEnum) error {
	query := "UPDATE disbursements SET status = $1, updated_at = NOW() WHERE id = $2"
	result, err := tx.Exec(query, status, disbursementID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrDisbursementNotFound
	}

	return nil
}
