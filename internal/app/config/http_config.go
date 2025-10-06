package config

import (
	"fmt"
	"strconv"
	"strings"
)

type HTTPConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func NewHTTPConfig(host string, port int) *HTTPConfig {
	return &HTTPConfig{
		Host: host,
		Port: port,
	}
}

func DefaultHTTPConfig() *HTTPConfig {
	return NewHTTPConfig(DefaultHTTPHost, DefaultHTTPPort)
}

func (hc *HTTPConfig) GetListenerAddr() string {
	return hc.Host + ":" + strconv.Itoa(hc.Port)
}

// flag.Value ==================================

func (hc *HTTPConfig) String() string {
	return fmt.Sprintf("%s:%v", hc.Host, hc.Port)
}

func (hc *HTTPConfig) Set(s string) error {
	params := strings.Split(s, ":")
	if len(params) < 2 {
		return fmt.Errorf("invalid http config format [%s], example: localhost:8080", s)
	}

	hc.Host = params[0]
	hc.Port, _ = strconv.Atoi(params[1])

	return nil
}

// =============================================
