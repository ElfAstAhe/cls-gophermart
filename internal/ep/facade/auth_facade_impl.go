package facade

import (
	"context"
	"encoding/json"
	"io"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/service"
	v1 "github.com/ElfAstAhe/cls-gophermart/internal/ep/dto/v1"
)

type AuthFacadeImpl struct {
	authService service.AuthService
	userService service.UserService
	log         logger.Logger
}

func NewAuthFacadeImpl(authService service.AuthService, userService service.UserService, logger logger.Logger) *AuthFacadeImpl {
	return &AuthFacadeImpl{
		authService: authService,
		userService: userService,
		log:         logger.GetLogger("AuthFacadeImpl"),
	}
}

func (a *AuthFacadeImpl) RegisterUser(ctx context.Context, register io.Reader) (string, error) {
	dec := json.NewDecoder(register)
	var dto v1.RegisterUserDto
	err := dec.Decode(&dto)
	if err != nil {
		return "", err
	}

	err = a.userService.Register(ctx, dto.Username, dto.Password)
	if err != nil {
		return "", err
	}

	jwtString, err := a.authService.Login(ctx, dto.Username, dto.Password)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}

func (a *AuthFacadeImpl) LoginUser(ctx context.Context, login io.Reader) (string, error) {
	dec := json.NewDecoder(login)
	var dto v1.LoginDto
	err := dec.Decode(&dto)
	if err != nil {
		return "", err
	}

	jwtString, err := a.authService.Login(ctx, dto.Username, dto.Password)
	if err != nil {
		return "", err
	}

	return jwtString, nil
}
