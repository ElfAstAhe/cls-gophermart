package repository

import (
	"context"
	"database/sql"
	"errors"

	_db "github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	_rep "github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
	"github.com/google/uuid"
)

const (
	findUserSql       string = `select id, username, password, disabled from users where id = $1`
	findUserByNameSql string = `select id, username, password, disabled from users where name = $1`
	createUserSql     string = `insert into users(id, username, password, disabled) values( $1, $2, $3, $4)`
	changeUserSql     string = `update users set password = $2, disabled = $3 where id = $1`
	softDeleteUserSql string = `update users set disabled = true where id = $1`
)

type UserRepositoryImpl struct {
	db          _db.DB
	accountRepo _rep.AccountRepository
}

func NewUserRepo(db _db.DB, accountRepo _rep.AccountRepository) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db:          db,
		accountRepo: accountRepo,
	}
}

func (u *UserRepositoryImpl) Find(ctx context.Context, id string) (*_mod.User, error) {
	row := u.db.GetDB().QueryRowContext(ctx, findUserSql, id)

	model := _mod.User{}
	err := row.Scan(&model.ID, &model.Username, &model.Password, &model.Disabled)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (u *UserRepositoryImpl) FindByName(ctx context.Context, username string) (*_mod.User, error) {
	row := u.db.GetDB().QueryRowContext(ctx, findUserByNameSql, username)

	model := _mod.User{}
	err := row.Scan(&model.ID, &model.Username, &model.Password, &model.Disabled)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (u *UserRepositoryImpl) Create(ctx context.Context, user *_mod.User) (*_mod.User, error) {
	if err := u.validateInstance(user); err != nil {
		return nil, err
	}
	if err := u.validateBL(ctx, user); err != nil {
		return nil, err
	}

	user.ID = uuid.New().String()

	_, err := u.db.GetDB().ExecContext(ctx, createUserSql, user.ID, user.Username, user.Password, user.Disabled)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepositoryImpl) Change(ctx context.Context, user *_mod.User) (*_mod.User, error) {
	if err := u.validateInstance(user); err != nil {
		return nil, err
	}

	_, err := u.db.GetDB().ExecContext(ctx, changeUserSql, user.ID, user.Password, user.Disabled)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	_, err := u.db.GetDB().ExecContext(ctx, softDeleteUserSql, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepositoryImpl) validateInstance(user *_mod.User) error {
	if user == nil {
		return _err.NewModelValidationError("user", "user is nil", nil)
	}

	if user.Username == "" {
		return _err.NewModelValidationError("user", "username is empty", nil)
	}

	if user.Password == "" {
		return _err.NewModelValidationError("user", "password is empty", nil)
	}

	return nil
}

func (u *UserRepositoryImpl) validateBL(ctx context.Context, user *_mod.User) error {
	model, err := u.FindByName(ctx, user.Username)
	if err != nil {
		return err
	}
	if model != nil {
		return _err.NewModelAlreadyExistsError("user")
	}

	return nil
}
