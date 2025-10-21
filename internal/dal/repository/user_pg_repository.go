package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/config/db"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
	"github.com/google/uuid"
)

const (
	pgFindUserSql       string = `select id, username, password, disabled from users where id = $1`
	pgFindUserByNameSql string = `select id, username, password, disabled from users where username = $1`
	pgCreateUserSql     string = `insert into users(id, username, password, disabled) values( $1, $2, $3, $4)`
	pgChangeUserSql     string = `update users set password = $2, disabled = $3 where id = $1`
	pgSoftDeleteUserSql string = `update users set disabled = true where id = $1`
)

type UserPgRepository struct {
	db          db.DB
	accountRepo repository.AccountRepository
}

func NewUserPgRepository(db db.DB, accountRepo repository.AccountRepository) *UserPgRepository {
	return &UserPgRepository{
		db:          db,
		accountRepo: accountRepo,
	}
}

func (u *UserPgRepository) Find(ctx context.Context, id string) (*model.User, error) {
	return u.findSingle(ctx, pgFindUserSql, id)
}

func (u *UserPgRepository) FindByName(ctx context.Context, username string) (*model.User, error) {
	return u.findSingle(ctx, pgFindUserByNameSql, username)
}

func (u *UserPgRepository) findSingle(ctx context.Context, query string, param any) (*model.User, error) {
	row := u.db.GetDB().QueryRowContext(ctx, query, param)

	user := model.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Disabled)

	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserPgRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
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

func (u *UserPgRepository) Change(ctx context.Context, user *model.User) (*model.User, error) {
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

func (u *UserPgRepository) validateInstance(user *model.User) error {
	if user == nil {
		return errs.NewModelValidationError("user", "user is nil", nil)
	}

	if user.Username == "" {
		return errs.NewModelValidationError("user", "username is empty", nil)
	}

	if user.Password == "" {
		return errs.NewModelValidationError("user", "password is empty", nil)
	}

	return nil
}

func (u *UserPgRepository) validateBL(ctx context.Context, user *model.User) error {
	founded, err := u.FindByName(ctx, user.Username)
	if err != nil {
		return err
	}
	if founded != nil {
		return errs.NewModelAlreadyExistsError("user", user.Username)
	}

	return nil
}
