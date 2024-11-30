package repository_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"simple-wallet-app/internal/sqlutil"
	"simple-wallet-app/module/wallet/entity"
	"simple-wallet-app/module/wallet/internal/repository"
	"simple-wallet-app/module/wallet/internal/testutil"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type userRepositorySuite struct {
	suite.Suite

	mockedDB  *sql.DB
	mockedSQL sqlmock.Sqlmock
	sut       *repository.UserRepository
}

func (s *userRepositorySuite) SetupTest() {
	mockDB, sqlMock, _ := sqlmock.New()
	s.mockedDB = mockDB
	s.mockedSQL = sqlMock

	s.sut = repository.NewUserRepository(s.mockedDB)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(userRepositorySuite))
}

func (s *userRepositorySuite) TestGetByID() {
	type output struct {
		expectedRes *entity.User
		expectedErr error
	}
	type input struct {
		param  *entity.User
		mockFn func(sm *sqlmock.ExpectedQuery, in input, out output)
		tx     sqlutil.DatabaseTransaction
	}

	defaultParam := testutil.DummyUser()
	defaultRes := testutil.DummyUser()

	tests := []struct {
		name string
		in   input
		out  output
	}{

		{
			name: "positive case: success",
			in: input{
				param: defaultParam,
				mockFn: func(sm *sqlmock.ExpectedQuery, in input, out output) {
					p := in.param
					r := out.expectedRes
					mockRows := s.mockedSQL.NewRows([]string{"id", "name", "balance", "pending_balance", "created_at", "updated_at"})
					mockRows.AddRow(r.ID, r.Name, r.Balance, r.PendingBalance, r.CreatedAt, r.UpdatedAt)
					sm.WithArgs(p.ID).WillReturnRows(mockRows)
				},
			},
			out: output{
				expectedRes: defaultRes,
				expectedErr: nil,
			},
		},
		{
			name: "negative case: query error",
			in: input{
				param: defaultParam,
				mockFn: func(sm *sqlmock.ExpectedQuery, in input, out output) {
					p := in.param
					sm.WithArgs(p.ID).WillReturnError(sql.ErrConnDone)
				},
			},
			out: output{
				expectedRes: nil,
				expectedErr: sql.ErrConnDone,
			},
		},
		{
			name: "negative case: empty row error",
			in: input{
				param: defaultParam,
				mockFn: func(sm *sqlmock.ExpectedQuery, in input, out output) {
					p := in.param
					mockRows := s.mockedSQL.NewRows([]string{"id", "name", "balance", "pending_balance", "created_at", "updated_at"})
					sm.WithArgs(p.ID).WillReturnRows(mockRows)
				},
			},
			out: output{
				expectedRes: nil,
				expectedErr: entity.ErrUserNotFound,
			},
		},
		{
			name: "negative case: scan error",
			in: input{
				param: defaultParam,
				mockFn: func(sm *sqlmock.ExpectedQuery, in input, out output) {
					p := in.param
					mockRows := s.mockedSQL.NewRows([]string{"id", "name", "balance", "pending_balance", "created_at", "updated_at"})
					mockRows.AddRow(p.ID, p.Name, "wrong balance", p.PendingBalance, p.CreatedAt, p.UpdatedAt)
					sm.WithArgs(p.ID).WillReturnRows(mockRows)
				},
			},
			out: output{
				expectedRes: nil,
				expectedErr: fmt.Errorf("sql: Scan error on column index 2, name \"balance\": can't convert wrong balance to decimal: exponent is not numeric"),
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// GIVEN
			query := regexp.QuoteMeta(`SELECT id, name, balance, pending_balance, created_at, updated_at FROM users WHERE id = $1`)
			if tt.in.mockFn != nil {
				mockSql := s.mockedSQL.ExpectQuery(query)
				tt.in.mockFn(mockSql, tt.in, tt.out)
			}

			// WHEN
			res, err := s.sut.GetByID(tt.in.tx, tt.in.param.ID)

			// THEN
			if (tt.out.expectedErr == nil && err != nil) || (tt.out.expectedErr != nil && !assert.EqualError(s.T(), err, tt.out.expectedErr.Error())) {
				s.T().Errorf("TestGetByID() error = %v, wantErr %v", err, tt.out.expectedErr)
			} else if tt.out.expectedRes != nil {
				s.Equal(tt.out.expectedRes, res)
				s.Nil(err)
			} else {
				s.EqualError(tt.out.expectedErr, err.Error())
			}
		})
	}
}

func (s *userRepositorySuite) TestHoldBalance() {
	type output struct {
		expectedErr error
	}
	type input struct {
		param  *entity.User
		amount decimal.Decimal
		mockFn func(sm *sqlmock.ExpectedQuery, in input, out output)
		tx     sqlutil.DatabaseTransaction
	}

	defaultParam := testutil.DummyUser()

	tests := []struct {
		name string
		in   input
		out  output
	}{
		{
			name: "positive case: hold balance success",
			in: input{
				param:  defaultParam,
				amount: decimal.NewFromInt(100),
				mockFn: func(sm *sqlmock.ExpectedQuery, in input, out output) {
					p := in.param
					mockRows := s.mockedSQL.NewRows([]string{"id", "name", "balance", "pending_balance", "created_at", "updated_at"})
					mockRows.AddRow(p.ID, p.Name, p.Balance, p.PendingBalance, p.CreatedAt, p.UpdatedAt)
					sm.WithArgs(p.ID, in.amount).WillReturnRows(mockRows)
				},
			},
			out: output{
				expectedErr: nil,
			},
		},
		{
			name: "negative case: hold balance insufficient funds",
			in: input{
				param:  defaultParam,
				amount: decimal.NewFromInt(1000),
				mockFn: func(sm *sqlmock.ExpectedQuery, in input, out output) {
					p := in.param
					sm.WithArgs(p.ID, in.amount).WillReturnError(fmt.Errorf("insufficient funds"))
				},
			},
			out: output{
				expectedErr: fmt.Errorf("insufficient funds"),
			},
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			// GIVEN
			query := regexp.QuoteMeta(`UPDATE users SET balance = balance - $1, pending_balance = pending_balance + $1 WHERE id = $2`)
			if tt.in.mockFn != nil {
				mockSql := s.mockedSQL.ExpectQuery(query)
				tt.in.mockFn(mockSql, tt.in, tt.out)
			}

			// WHEN
			err := s.sut.HoldBalance(tt.in.tx, tt.in.param.ID, tt.in.amount)

			// THEN
			if (tt.out.expectedErr == nil && err != nil) || (tt.out.expectedErr != nil && !assert.EqualError(s.T(), err, tt.out.expectedErr.Error())) {
				s.T().Errorf("%s error = %v, wantErr %v", tt.name, err, tt.out.expectedErr)
			}
		})
	}
}
