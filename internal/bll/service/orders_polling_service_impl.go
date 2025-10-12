package service

import (
	"context"
	"sync"
	"time"

	_log "github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	_repo "github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
)

type OrdersPollingServiceImpl struct {
	parentCtx  context.Context
	stopCtx    context.Context
	cancelFunc context.CancelFunc
	queue      chan string
	scheduler  *time.Timer
	wg         *sync.WaitGroup
	orderRepo  _repo.OrderRepository
	log        _log.AppLogger
	//    lsClient   _hc.LSClient
}

func NewOrdersPollingService(ctx context.Context, orderRepo _repo.OrderRepository, logger _log.AppLogger) OrdersPollingService {
	return &OrdersPollingServiceImpl{
		parentCtx:  ctx,
		stopCtx:    nil,
		cancelFunc: nil,
		queue:      nil,
		scheduler:  nil,
		wg:         &sync.WaitGroup{},
		orderRepo:  orderRepo,
		log:        logger.GetLogger("OrdersPollingService"),
	}
}

func (ops *OrdersPollingServiceImpl) Start(ctx context.Context) error {

	// stop ctx
	ops.stopCtx, ops.cancelFunc = context.WithCancel(ops.parentCtx)
	// timer
	if ops.scheduler == nil {
		ops.scheduler = time.NewTimer(10 * time.Second)
	} else {
		ops.scheduler.Reset(10 * time.Second)
	}
	// queue
	ops.queue = make(chan string, 128)
	// launch timer event method
	go ops.timerEventListener(ops.scheduler, ops.stopCtx)
	// launch workers
	// ToDo: implement, create poll workers
	// ..

	return nil
}

func (ops *OrdersPollingServiceImpl) Stop(ctx context.Context) error {
	// stop timer
	if ops.scheduler != nil {
		ops.scheduler.Stop()
	}
	// stop timer event and workers
	if ops.cancelFunc != nil {
		ops.cancelFunc()
	}
	if ops.queue != nil {
		close(ops.queue)
	}

	ops.wg.Wait()

	return nil
}

func (ops *OrdersPollingServiceImpl) timerEventListener(timer *time.Timer, stopCtx context.Context) {

	for {
		select {
		case event <- timer.C:
			processEvent(event)

		case <-stopCtx.Done():
			return
		}
	}
}
