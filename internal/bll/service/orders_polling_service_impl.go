package service

import (
	"context"
	"sync"
	"time"

	_log "github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	_repo "github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	_lsc "github.com/ElfAstAhe/cls-gophermart/pkg/client/loyalty"
)

type OrdersPollingServiceImpl struct {
	parentCtx      context.Context
	stopCtx        context.Context
	cancelFunc     context.CancelFunc
	queue          chan string
	schedulerTimer *time.Timer
	wg             *sync.WaitGroup
	baseURI        string
	orderRepo      _repo.OrderRepository
	log            _log.AppLogger
	lsClient       _lsc.LSClient
}

func NewOrdersPollingService(ctx context.Context, baseURI string, orderRepo _repo.OrderRepository, logger _log.AppLogger) OrdersPollingService {
	return &OrdersPollingServiceImpl{
		parentCtx:      ctx,
		stopCtx:        nil,
		cancelFunc:     nil,
		queue:          nil,
		schedulerTimer: nil,
		wg:             &sync.WaitGroup{},
		baseURI:        baseURI,
		orderRepo:      orderRepo,
		log:            logger.GetLogger("OrdersPollingService"),
	}
}

func (ops *OrdersPollingServiceImpl) Start(ctx context.Context) error {

	// stop ctx
	ops.stopCtx, ops.cancelFunc = context.WithCancel(ops.parentCtx)
	// timer
	if ops.schedulerTimer == nil {
		ops.schedulerTimer = time.NewTimer(10 * time.Second)
	} else {
		ops.schedulerTimer.Reset(10 * time.Second)
	}
	// queue
	ops.queue = make(chan string, 128)
	// launch timer event method
	go ops.timerEventListener(ops.schedulerTimer, ops.stopCtx)
	// launch workers
	// ToDo: implement, create poll workers
	// ..

	return nil
}

func (ops *OrdersPollingServiceImpl) Stop(ctx context.Context) error {
	// stop timer
	if ops.schedulerTimer != nil {
		ops.schedulerTimer.Stop()
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
		case <-timer.C:
			{
				ops.processEvent(stopCtx)
				timer.Reset(10 * time.Second)
			}
		case <-stopCtx.Done():
			return
		}
	}
}

func (ops *OrdersPollingServiceImpl) processEvent(stopCtx context.Context) {
	ops.log.Debug("order polling service scheduler timer event started")
	defer ops.log.Debug("order polling service scheduler timer event finished")
	queueLength := len(ops.queue)
	if queueLength != 0 {
		ops.log.Infof("order polling service scheduler timer event queue length, exit: [%v]", queueLength)

		return
	}
	orders, err := ops.orderRepo.ListUnfinished(stopCtx, _mod.OrderStatusNew, _mod.OrderStatusProcessing)
	if err != nil {
		ops.log.Errorf("error get unfinished orders: %v", err)
	}
	for _, order := range orders {
		select {
		case <-stopCtx.Done():
			return
		default:
			ops.queue <- order.ID

		}
	}
}
