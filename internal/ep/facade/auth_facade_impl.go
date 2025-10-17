package facade

import (
	"context"
	"io"

	_log "github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	_svc "github.com/ElfAstAhe/cls-gophermart/internal/bll/service"
)

type AuthFacadeImpl struct {
	authService _svc.AuthService
	log         _log.AppLogger
}

func NewAuthFacadeImpl(authService _svc.AuthService, logger _log.AppLogger) *AuthFacadeImpl {
	return &AuthFacadeImpl{
		authService: authService,
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
