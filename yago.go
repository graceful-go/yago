package yago

import (
	"context"
	"fmt"
	"net/http"
)

// YagoConfig
type YagoConfig struct {
	Port    uint32 `json:"port" default:"8080"`
	Timeout uint32 `json:"timeout" default:"1000"`
}

type YagoHandler interface {
	Handler() http.Handler
	Pattern() string
	Type() string
}

type Yago struct {
	cfg      *YagoConfig
	handlers []YagoHandler
	logger   Logger
}

func New(opts ...Option) (*Yago, error) {

	logger := &DefaultLogger{}

	y := &Yago{logger: logger}
	for _, opt := range opts {
		opt(y)
	}

	if err := y.check(); err != nil {
		return nil, err
	}

	return y, nil
}

func (y *Yago) check() error {
	return nil
}

// Start will block current process and start up a http server
func (y *Yago) Start(ctx context.Context) error {

	for _, handler := range y.handlers {
		http.Handle(handler.Pattern(), handler.Handler())
		y.logger.Log("[YagoServer] Start server for: ", handler.Pattern(), "[Type]", handler.Type())
	}

	y.logger.Log("[YagoServer] Server Startup at", y.cfg.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", y.cfg.Port), nil)
}
