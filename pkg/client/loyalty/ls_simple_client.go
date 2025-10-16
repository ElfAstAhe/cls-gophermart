package loyalty

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"time"

	_log "github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	_dto "github.com/ElfAstAhe/cls-gophermart/pkg/client/loyalty/dto"
)

const (
	lsGetOrderPathParam string = "{orderNumber}"
)

type LSSimpleClient struct {
	client  *http.Client
	baseURL string
	log     _log.AppLogger
}

func NewLSSimpleClient(baseURL string, timeOut time.Duration, logger _log.AppLogger) (*LSSimpleClient, error) {
	return &LSSimpleClient{
		client: &http.Client{
			Timeout: timeOut,
		},
		baseURL: baseURL,
		log:     logger.GetLogger("ls-simple-client"),
	}, nil
}

func (lc *LSSimpleClient) GetOrder(ctx context.Context, orderNumber string, token string) (*_dto.LSOrderDto, error) {
	getOrderUrl := path.Join(lc.baseURL, lsGetOrderPathParam)
	req, err := http.NewRequest(http.MethodGet, getOrderUrl, nil)
	if err != nil {
		return nil, NewLSClientError(fmt.Sprintf("error creating request to get order [%s]", orderNumber), -1, err)
	}
	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, NewLSClientError(fmt.Sprintf("error executing request to get order [%s]", orderNumber), -1, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, NewLSClientError(fmt.Sprintf("error executing request to get order [%s] with status code [%v]", orderNumber, resp.StatusCode), resp.StatusCode, nil)
	}
	decoder := json.NewDecoder(resp.Body)
	dto := &_dto.LSOrderDto{}
	if err := decoder.Decode(dto); err != nil {
		return nil, NewLSClientError(fmt.Sprintf("error executing request to get order [%s]", orderNumber), -1, err)
	}

	return dto, nil
}
