package config

import (
	"fmt"
	"strconv"
	"strings"
)

type HTTPConfig struct {
	Schema string `json:"schema"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
}

func NewHTTPConfig(schema string, host string, port int) *HTTPConfig {
	return &HTTPConfig{
		Schema: schema,
		Host:   host,
		Port:   port,
	}
}

func DefaultHTTPConfig() *HTTPConfig {
	return NewHTTPConfig(DefaultHTTPSchema, DefaultHTTPHost, DefaultHTTPPort)
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
