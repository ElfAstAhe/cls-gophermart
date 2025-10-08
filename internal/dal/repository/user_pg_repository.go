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
	pgFindUserSql       string = `select id, username, password, disabled from users where id = $1`
	pgFindUserByNameSql string = `select id, username, password, disabled from users where name = $1`
	pgCreateUserSql     string = `insert into users(id, username, password, disabled) values( $1, $2, $3, $4)`
	pgChangeUserSql     string = `update users set password = $2, disabled = $3 where id = $1`
	pgSoftDeleteUserSql string = `update users set disabled = true where id = $1`
)

type UserPgRepository struct {
	db          _db.DB
	accountRepo _rep.AccountRepository
}

func NewUserPgRepository(db _db.DB, accountRepo _rep.AccountRepository) *UserPgRepository {
	return &UserPgRepository{
		db:          db,
		accountRepo: accountRepo,
	}
}

func (u *UserPgRepository) Find(ctx context.Context, id string) (*_mod.User, error) {
	return u.findSingle(ctx, pgFindUserSql, id)
}

func (u *UserPgRepository) FindByName(ctx context.Context, username string) (*_mod.User, error) {
	return u.findSingle(ctx, pgFindUserByNameSql, username)
}

func (u *UserPgRepository) findSingle(ctx context.Context, query string, param any) (*_mod.User, error) {
	row := u.db.GetDB().QueryRowContext(ctx, query, param)

	model := _mod.User{}
	err := row.Scan(&model.ID, &model.Username, &model.Password, &model.Disabled)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &model, nil
}

func (u *UserPgRepository) Create(ctx context.Context, user *_mod.User) (*_mod.User, error) {
	if err := u.validateInstance(user); err != nil {
		return nil, err
	}
	if err := u.validateBL(ctx, user); err != nil {
		return nil, err
	}

	user.ID = uuid.New().String()

	_, err := u.db.GetDB().ExecContext(ctx, pgCreateUserSql, user.ID, user.Username, user.Password, user.Disabled)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserPgRepository) Change(ctx context.Context, user *_mod.User) (*_mod.User, error) {
	if err := u.validateInstance(user); err != nil {
		return nil, err
	}

	_, err := u.db.GetDB().ExecContext(ctx, pgChangeUserSql, user.ID, user.Password, user.Disabled)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserPgRepository) SoftDelete(ctx context.Context, id string) error {
	_, err := u.db.GetDB().ExecContext(ctx, pgSoftDeleteUserSql, id)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserPgRepository) validateInstance(user *_mod.User) error {
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

func (u *UserPgRepository) validateBL(ctx context.Context, user *_mod.User) error {
	model, err := u.FindByName(ctx, user.Username)
	if err != nil {
		return err
	}
	if model != nil {
		return _err.NewModelAlreadyExistsError("user", user.Username)
	}

	return nil
}
