package repository

import (
	"context"

	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type AccountRepository interface {
	// Find search by synthetic key
	Find(ctx context.Context, id string) (*model.Account, error)
	// FindByUser search by bl key (one to one link)
	FindByUser(ctx context.Context, userID string) (*model.Account, error)
	// Create new account
	Create(ctx context.Context, userID string, account *model.Account) (*model.Account, error)
	// Change account attributes
	Change(ctx context.Context, userID string, account *model.Account) (*model.Account, error)
	// GetBalance get account current full balance info
	GetBalance(ctx context.Context, id string) (*model.AccountBalance, error)
}
