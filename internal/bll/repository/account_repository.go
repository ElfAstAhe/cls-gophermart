package repository

import (
	"context"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type AccountRepository interface {
	// Find search by synthetic key
	Find(ctx context.Context, ID string) (*_mod.Account, error)
	// FindByUser search by bl key (one to one link)
	FindByUser(ctx context.Context, userID string) (*_mod.Account, error)
	// Create new account
	Create(ctx context.Context, userID string, account *_mod.Account) (*_mod.Account, error)
	// Change account attributes
	Change(ctx context.Context, userID string, account *_mod.Account) (*_mod.Account, error)
	// SoftDelete set deleted flag
	SoftDelete(ctx context.Context, ID string) error
	// Remove physical remove account
	Remove(ctx context.Context, ID string) error
}
