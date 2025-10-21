package facade

import (
	"context"
	"io"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/service"
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
	//TODO implement me
	panic("implement me")
}

func (a *AuthFacadeImpl) LoginUser(ctx context.Context, login io.Reader) (string, error) {
	//TODO implement me
	panic("implement me")
}
