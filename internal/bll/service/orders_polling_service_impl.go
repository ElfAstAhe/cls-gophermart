package service

import (
	"context"
	"errors"
	"sync"
	"time"

	_log "github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	_mod "github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	_repo "github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	_lsc "github.com/ElfAstAhe/cls-gophermart/pkg/client/loyalty"
	_dto "github.com/ElfAstAhe/cls-gophermart/pkg/client/loyalty/dto"
	_err "github.com/ElfAstAhe/cls-gophermart/pkg/error"
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
	ops.log.Info("Starting Order Polling Service")
	defer ops.log.Info("Done starting Order Polling Service")
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
	go ops.timerEventListener(ops.stopCtx, ops.schedulerTimer)
	// launch workers
	if err := ops.createWorkerPool(ops.stopCtx, 4); err != nil {
		return err
	}

	return nil
}

func (ops *OrdersPollingServiceImpl) Stop(ctx context.Context) error {
	ops.log.Info("Stopping Order Polling Service")
	defer ops.log.Info("Done stopping Order Polling Service")
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

func (ops *OrdersPollingServiceImpl) timerEventListener(stopCtx context.Context, timer *time.Timer) {
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
	ops.log.Infof("processing order polling service scheduler timer event queue length: %d", queueLength)
	if queueLength != 0 {
		ops.log.Info("processing order polling service scheduler timer event queue not empty, exit")

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

func (ops *OrdersPollingServiceImpl) createWorkerPool(stopCtx context.Context, workerCount int) error {
	if workerCount <= 0 {
		return _err.NewAppInvalidArgumentError("workerCount", workerCount)
	}

	for index := 0; index < workerCount; index++ {
		ops.wg.Add(1)
		go ops.worker(stopCtx, index, 1*time.Second)
	}

	return nil
}

func (ops *OrdersPollingServiceImpl) worker(stopCtx context.Context, index int, sleepTime time.Duration) {
	ops.log.Infof("order polling service worker index [%v] started", index)
	defer ops.log.Infof("order polling service worker index [%v] finished", index)
	defer ops.wg.Done()

	// infinite circle
	for {
		ops.log.Infof("order polling service worker index [%v] start poll iteration", index)
		select {

		case <-stopCtx.Done():
			return
		case id := <-ops.queue:
			ops.log.Infof("processing order id=[%s] worker index [%v] start", id, index)
			if err := ops.processSingleOrder(id); err != nil {
				var lsErr *_lsc.LSClientError
				if errors.As(err, &lsErr) {
					go ops.processLSError(lsErr, id)
				} else {
					ops.log.Errorf("error processing order id=[%s] worker index [%v] error: %v", id, index, err)
				}
			}

			ops.log.Infof("processing order id=[%s] worker index [%v] finished, sleep [%v]", id, index, sleepTime)
			time.Sleep(sleepTime)
		}

		ops.log.Infof("order polling service worker index [%v] finish poll iteration", index)
	}
}

func (ops *OrdersPollingServiceImpl) processSingleOrder(orderID string) error {
	dto, err := ops.lsClient.GetOrder(orderID)
	// check client response
	if err != nil {
		return err
	}
	newModel, err := ops.toModel(dto)
	if err != nil {
		return err
	}
	newModel.ID = orderID

	oldModel, err := ops.orderRepo.Find(ops.stopCtx, orderID)
	if err != nil {
		return err
	}

	// business logic
	if oldModel.Status != newModel.Status {
		if _, err := ops.orderRepo.Change(ops.stopCtx, newModel); err != nil {
			return err
		}
		ops.log.Infof("order id=[%s] status changed to [%s] from [%s], changes stored", orderID, newModel.Status, oldModel.Status)
	} else {
		ops.log.Infof("order id=[%s] no changes", orderID)
	}

	return nil
}

func (ops *OrdersPollingServiceImpl) processLSError(lsErr *_lsc.LSClientError, id string) {

}

func (ops *OrdersPollingServiceImpl) toModel(dto *_dto.LSOrderDto) (*_mod.Order, error) {
	if dto == nil {
		return nil, _err.NewAppInvalidArgumentError("order dto", nil)
	}
	status, err := ops.toModelStatus(dto.Status)
	if err != nil {
		return nil, err
	}

	return _mod.NewOrder(dto.Number, status, dto.AccrualAmount, time.Now()), nil
}

func (ops *OrdersPollingServiceImpl) toModelStatus(dtoStatus string) (string, error) {
	switch dtoStatus {
	case _dto.LSOrderStatusRegistered:
		return _mod.OrderStatusNew, nil
	case _dto.LSOrderStatusProcessing:
		return _mod.OrderStatusProcessing, nil
	case _dto.LSOrderStatusInvalid:
		return _mod.OrderStatusInvalid, nil
	case _dto.LSOrderStatusProcessed:
		return _mod.OrderStatusProcessed, nil
	}

	return "", _err.NewAppInvalidArgumentError("dtoStatus", dtoStatus)
}
