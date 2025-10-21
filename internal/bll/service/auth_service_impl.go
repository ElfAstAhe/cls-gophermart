package service

import (
	"context"
	"strings"

	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/service/auth"
	errs "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

type AuthServiceImpl struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		userRepo: userRepo,
	}
}

func (as *AuthServiceImpl) Login(ctx context.Context, username string, password string) (string, error) {
	if err := as.validateArgs(username, password); err != nil {
		return "", err
	}

	user, err := as.userRepo.FindByName(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errs.NewModelNotExistsError("user", username)
	}

	if err := as.checkPassword(user, password); err != nil {
		return "", err
	}

	jwtString, err := auth.NewJWTString(false, user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}

func (as *AuthServiceImpl) validateArgs(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return errs.NewAppInvalidArgumentError("username", username)
	}
	if strings.TrimSpace(password) == "" {
		return errs.NewAppInvalidArgumentError("password", password)
	}

	return nil
}

func (as *AuthServiceImpl) checkPassword(user *model.User, password string) error {
	if password != user.Password {
		return errs.NewAuthPasswordIncorrectError(user.Username)
	}

	return nil
}
