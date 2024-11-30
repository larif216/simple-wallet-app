package testutil

import (
	"simple-wallet-app/module/wallet/entity"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func Name(tags []string, testcaseName string) string {
	return strings.Join(tags, ",") + ": " + testcaseName
}

func DummyUser() *entity.User {
	return &entity.User{
		ID:             1,
		Name:           "Lutfi",
		Balance:        decimal.New(10000, 0),
		PendingBalance: decimal.New(0, 0),
		CreatedAt:      time.Date(2023, 01, 22, 10, 10, 10, 0, time.UTC),
		UpdatedAt:      time.Date(2023, 01, 22, 10, 10, 10, 0, time.UTC),
	}
}
