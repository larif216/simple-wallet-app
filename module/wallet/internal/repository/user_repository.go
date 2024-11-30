package repository

import (
	"database/sql"
	"simple-wallet-app/internal/sqlutil"
	"simple-wallet-app/module/wallet/entity"
	"time"

	"github.com/shopspring/decimal"
)

type userRecord struct {
	ID             uint32          `db:"id"`
	Name           string          `db:"name"`
	Balance        decimal.Decimal `db:"balance"`
	PendingBalance decimal.Decimal `db:"pending_balance"`
	CreatedAt      time.Time       `db:"created_at"`
	UpdatedAt      time.Time       `db:"updated_at"`
}

func (c *userRecord) toEntity() *entity.User {
	return &entity.User{
		ID:             c.ID,
		Name:           c.Name,
		Balance:        c.Balance,
		PendingBalance: c.PendingBalance,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetByID(tx sqlutil.DatabaseTransaction, userID uint32) (*entity.User, error) {
	var userRecord userRecord

	query := "SELECT id, name, balance, pending_balance, created_at, updated_at FROM users WHERE id = $1"
	err := tx.QueryRow(
		query,
		userID,
	).Scan(
		&userRecord.ID,
		&userRecord.Name,
		&userRecord.Balance,
		&userRecord.PendingBalance,
		&userRecord.CreatedAt,
		&userRecord.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrUserNotFound
		}
		return nil, err
	}

	return userRecord.toEntity(), nil
}

func (r *UserRepository) HoldBalance(tx sqlutil.DatabaseTransaction, userID uint32, amount decimal.Decimal) error {
	query := "UPDATE users SET balance = balance - $1, pending_balance = pending_balance + $1 WHERE id = $2"
	result, err := tx.Exec(query, amount, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) ReleaseBalance(tx sqlutil.DatabaseTransaction, userID uint32, amount decimal.Decimal, withReversal bool) error {
	var query string
	if withReversal {
		query = "UPDATE users SET balance = balance + $1, pending_balance = pending_balance - $1 WHERE id = $2"
	} else {
		query = "UPDATE users SET pending_balance = pending_balance - $1 WHERE id = $2"
	}

	result, err := tx.Exec(query, amount, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return entity.ErrUserNotFound
	}

	return nil
}
