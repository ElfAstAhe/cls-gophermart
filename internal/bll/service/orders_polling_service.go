package service

import (
	"context"
)

type OrdersPollingService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
