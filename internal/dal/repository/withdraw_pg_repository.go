package repository

import (
	"context"

	_db "github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type WithdrawRepositoryImpl struct {
	db _db.DB
}

func (w WithdrawRepositoryImpl) Find(ctx context.Context, id string) (*_mod.Withdraw, error) {
	//TODO implement me
	panic("implement me")
}

func (w WithdrawRepositoryImpl) ListByUser(ctx context.Context, userId string) ([]*_mod.Withdraw, error) {
	//TODO implement me
	panic("implement me")
}

func (w WithdrawRepositoryImpl) WithdrawsByUser(ctx context.Context, userId string) (float64, error) {
	//TODO implement me
	panic("implement me")
}
