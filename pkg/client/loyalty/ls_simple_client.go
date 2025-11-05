package loyalty

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ElfAstAhe/cls-gophermart/internal/app/logger"
	"github.com/ElfAstAhe/cls-gophermart/pkg/client/loyalty/dto"
)

const (
	lsGetOrderPathParam string = "{orderNumber}"
)

type LSSimpleClient struct {
	client  *http.Client
	baseURL string
	log     logger.Logger
}

func NewLSSimpleClient(baseURL string, timeOut time.Duration, logger logger.Logger) *LSSimpleClient {
	return &LSSimpleClient{
		client: &http.Client{
			Timeout: timeOut,
		},
		baseURL: baseURL,
		log:     logger.GetLogger("ls-simple-client"),
	}
}

func (lc *LSSimpleClient) GetOrder(ctx context.Context, orderNumber string) (*dto.LSOrderDto, error) {
	getOrderUrl, err := url.Parse(lc.baseURL)
	if err != nil {
		return nil, err
	}
	getOrderUrl = getOrderUrl.JoinPath(orderNumber)

	req, err := http.NewRequest(http.MethodGet, getOrderUrl.String(), nil)
	if err != nil {
		return nil, NewLSClientError(fmt.Sprintf("error creating request to get order [%s]", orderNumber), -1, err)
	}
	resp, err := lc.client.Do(req)
	if err != nil {
		return nil, NewLSClientError(fmt.Sprintf("error executing request to get order [%s]", orderNumber), -1, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, NewLSClientError(fmt.Sprintf("error executing request to get order [%s] with status code [%v]", orderNumber, resp.StatusCode), resp.StatusCode, nil)
	}
	decoder := json.NewDecoder(resp.Body)
	lsOrder := &dto.LSOrderDto{}
	if err := decoder.Decode(lsOrder); err != nil {
		return nil, NewLSClientError(fmt.Sprintf("error executing request to get order [%s]", orderNumber), -1, err)
	}

	return lsOrder, nil
}
