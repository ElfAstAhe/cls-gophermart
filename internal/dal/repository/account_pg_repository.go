package repository

import (
	"context"

	_db "github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type AccountRepositoryImpl struct {
	db _db.DB
}

func (a AccountRepositoryImpl) Find(ctx context.Context, id string) (*_mod.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a AccountRepositoryImpl) FindByUser(ctx context.Context, userID string) (*_mod.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a AccountRepositoryImpl) Create(ctx context.Context, userID string, account *_mod.Account) (*_mod.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (a AccountRepositoryImpl) Change(ctx context.Context, userID string, account *_mod.Account) (*_mod.Account, error) {
	//TODO implement me
	panic("implement me")
}
