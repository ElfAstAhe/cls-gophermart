package repository

import (
	"context"

	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type UserRepository interface {
	Find(ctx context.Context, id string) (*model.User, error)
	FindByName(ctx context.Context, username string) (*model.User, error)
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Change(ctx context.Context, user *model.User) (*model.User, error)
	SoftDelete(ctx context.Context, id string) error
}
