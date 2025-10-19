package service

import (
	"context"
	"strings"

	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	errors "github.com/ElfAstAhe/cls-gophermart/pkg/error"
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
		return "", errors.NewModelNotExistsError("user", username)
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
		return "", errors.NewModelAlreadyExistsError("user", username)
	}

	user = model.NewUser(username, password)
	user, err = as.userRepo.Create(ctx, user)
	if err != nil {
		return "", err
	}

	return as.Login(ctx, username, password)
}

func (as *AuthServiceImpl) validateArgs(username, password string) error {
	if strings.TrimSpace(username) == "" {
		return errors.NewAppInvalidArgumentError("username", username)
	}
	if strings.TrimSpace(password) == "" {
		return errors.NewAppInvalidArgumentError("password", password)
	}

	return nil
}

func (as *AuthServiceImpl) checkPassword(user *model.User, password string) error {
	if password != user.Password {
		return errors.NewAuthPasswordIncorrectError(user.Username)
	}

	return nil
}
