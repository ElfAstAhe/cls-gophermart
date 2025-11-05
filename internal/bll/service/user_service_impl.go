package service

import (
	"context"
	"strings"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

type UserServiceImpl struct {
	userRepo    repository.UserRepository
	accountRepo repository.AccountRepository
	log         logger.Logger
}

func NewUserServiceImpl(userRepo repository.UserRepository, accountRepo repository.AccountRepository, logger logger.Logger) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:    userRepo,
		accountRepo: accountRepo,
		log:         logger.GetLogger("UserServiceImpl"),
	}
}

func (us *UserServiceImpl) Register(ctx context.Context, username, password string) error {
	if err := us.validateArgs(username, password); err != nil {
		return err
	}

	user, err := us.userRepo.FindByName(ctx, username)
	if err != nil {
		return err
	}
	if user != nil {
		return errs.NewModelAlreadyExistsError("user", username)
	}

	user, err = us.userRepo.Create(ctx, model.NewUser(username, password, false))
	if err != nil {
		return err
	}

	_, err = us.accountRepo.Create(ctx, user.ID, model.NewAccount("", "unknown"))
	if err != nil {
		return err
	}

	return nil
}

func (us *UserServiceImpl) validateArgs(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return errs.NewAppInvalidArgumentError("username", username)
	}
	if strings.TrimSpace(password) == "" {
		return errs.NewAppInvalidArgumentError("password", password)
	}

	return nil
}
