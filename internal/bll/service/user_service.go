package service

import "context"

type UserService interface {
	Register(ctx context.Context, username, password string) error
}
