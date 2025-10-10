package service

import (
	"context"
)

type AuthService interface {
	Login(ctx context.Context, username string, password string) (string, error)
	Register(ctx context.Context, username string, password string) (string, error)
}
