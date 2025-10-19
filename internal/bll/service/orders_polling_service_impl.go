package service

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/model"
	"github.com/ElfAstAhe/cls-gophermart/internal/bll/repository"
	"github.com/ElfAstAhe/cls-gophermart/pkg/client/loyalty"
	ldto "github.com/ElfAstAhe/cls-gophermart/pkg/client/loyalty/dto"
	errors "github.com/ElfAstAhe/cls-gophermart/pkg/error"
)

type OrdersPollingServiceImpl struct {
	parentCtx      context.Context
	stopCtx        context.Context
	cancelFunc     context.CancelFunc
	queue          chan string
	schedulerTimer *time.Timer
	wg             *sync.WaitGroup
	baseURI        string
	orderRepo      repository.OrderRepository
	log            logger.Logger
	lsClient       loyalty.LSClient
}

func NewOrdersPollingService(ctx context.Context, baseURI string, orderRepo repository.OrderRepository, logger logger.Logger) OrdersPollingService {
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

func (ops *OrdersPollingServiceImpl) Start() error {
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

func (ops *OrdersPollingServiceImpl) Stop() error {
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
	ops.log.Info("order polling service scheduler timer event started")
	defer ops.log.Info("order polling service scheduler timer event finished")
	queueLength := len(ops.queue)
	ops.log.Infof("processing order polling service scheduler timer event queue length: %d", queueLength)
	if queueLength != 0 {
		ops.log.Info("processing order polling service scheduler timer event queue not empty, exit")

		return
	}
	orders, err := ops.orderRepo.ListUnfinished(stopCtx, model.OrderStatusNew, model.OrderStatusProcessing)
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
		return errors.NewAppInvalidArgumentError("workerCount", workerCount)
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
				ops.log.Errorf("error processing order id=[%s] worker index [%v] error: %v", id, index, err)
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
	if dto == nil {
		ops.log.Warnf("order id [%s] not exists in loyalty system", orderID)
		return nil
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

func (ops *OrdersPollingServiceImpl) toModel(dto *ldto.LSOrderDto) (*model.Order, error) {
	if dto == nil {
		return nil, errors.NewAppInvalidArgumentError("order dto", nil)
	}
	status, err := ops.toModelStatus(dto.Status)
	if err != nil {
		return nil, err
	}

	return model.NewOrder(dto.Number, status, dto.AccrualAmount, time.Now()), nil
}

func (ops *OrdersPollingServiceImpl) toModelStatus(dtoStatus string) (string, error) {
	switch dtoStatus {
	case ldto.LSOrderStatusRegistered:
		return model.OrderStatusNew, nil
	case ldto.LSOrderStatusProcessing:
		return model.OrderStatusProcessing, nil
	case ldto.LSOrderStatusInvalid:
		return model.OrderStatusInvalid, nil
	case ldto.LSOrderStatusProcessed:
		return model.OrderStatusProcessed, nil
	}

	return "", errors.NewAppInvalidArgumentError("dtoStatus", dtoStatus)
}

func (ops *OrdersPollingServiceImpl) Add(ID string) {
	if strings.TrimSpace(ID) == "" {
		return
	}

	go func() {
		ops.queue <- ID
	}()
}
