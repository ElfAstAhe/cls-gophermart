package service

type OrdersPollingService interface {
	Start() error
	Stop() error
}
