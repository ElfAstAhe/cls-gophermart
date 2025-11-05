package facade

import (
	"context"
	"io"
)

type AuthFacade interface {
	RegisterUser(ctx context.Context, register io.Reader) (string, error)
	LoginUser(ctx context.Context, login io.Reader) (string, error)
}
