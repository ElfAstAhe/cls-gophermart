package service

import (
	"context"
	"strings"

	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
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
		return "", err
	}
	if user == nil {
		return "", _err.NewModelNotExistsError("user", username)
	}

	if err := as.checkPassword(user, password); err != nil {
		return "", err
	}

	// ToDo: generate jwt
	var jwtString string = "test.jwt.token"

	return jwtString, nil
}

func (as *AuthServiceImpl) Register(ctx context.Context, username string, password string) (string, error) {
	if err := as.validateArgs(username, password); err != nil {
		return "", err
	}

	user, err := as.userRepo.FindByName(ctx, username)
	if err != nil {
		return "", err
	}
	if user != nil {
		return "", _err.NewModelAlreadyExistsError("user", username)
	}

	user = _mod.NewUser(username, password)
	user, err = as.userRepo.Create(ctx, user)
	if err != nil {
		return "", err
	}

	return as.Login(ctx, username, password)
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

func (as *AuthServiceImpl) checkPassword(user *_mod.User, password string) error {
	if password != user.Password {
		return _err.NewAuthPasswordIncorrectError(user.Username)
	}

	return nil
}
