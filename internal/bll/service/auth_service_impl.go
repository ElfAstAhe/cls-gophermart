package service

import (
	"context"
	"strings"

	_repo "github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

type AuthServiceImpl struct {
	userRepo _repo.UserRepository
}

func NewAuthService(userRepo _repo.UserRepository) *AuthServiceImpl {
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

	}
}

func (as *AuthServiceImpl) Register(ctx context.Context, username string, password string) (string, error) {
	if err := as.validateArgs(username, password); err != nil {
		return "", err
	}

}

func (as *AuthServiceImpl) validateArgs(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return _err.NewAppInvalidArgumentError("username", username)
	}
	if strings.TrimSpace(password) == "" {
		return _err.NewAppInvalidArgumentError("password", password)
	}

	return nil
}
