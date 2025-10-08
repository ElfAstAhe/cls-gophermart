package repository

import (
	"context"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type UserRepository interface {
	Find(ctx context.Context, id string) (*_mod.User, error)
	FindByName(ctx context.Context, username string) (*_mod.User, error)
	Create(ctx context.Context, user *_mod.User) (*_mod.User, error)
	Change(ctx context.Context, user *_mod.User) (*_mod.User, error)
	SoftDelete(ctx context.Context, id string) error
}
