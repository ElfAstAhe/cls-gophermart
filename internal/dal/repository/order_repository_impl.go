package repository

import (
	"context"

	_db "github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
)

type OrderRepositoryImpl struct {
	db _db.DB
}

func (o OrderRepositoryImpl) Find(ctx context.Context, id string) (*_mod.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderRepositoryImpl) FindByNumber(ctx context.Context, number string) (*_mod.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderRepositoryImpl) ListByAccount(ctx context.Context, accountId string) ([]*_mod.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderRepositoryImpl) Create(ctx context.Context, accountId string, order *_mod.Order) (*_mod.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (o OrderRepositoryImpl) Change(ctx context.Context, accountId string, order *_mod.Order) (*_mod.Order, error) {
	//TODO implement me
	panic("implement me")
}
