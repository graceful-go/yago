package yago

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
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
	paths    map[string]YagoHandler
	mu       sync.RWMutex
}

func New(opts ...Option) (*Yago, error) {

	logger := &DefaultLogger{}

	y := &Yago{
		logger: logger,
		paths:  make(map[string]YagoHandler),
	}
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
	y.logger.Log("[YagoServer] Server Startup on port", y.cfg.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", y.cfg.Port), y)
}

func (y *Yago) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, handlerName := y.findHandler(r.URL.Path)
	if handler == nil {
		y.logger.Loglnf("[Yago] handler not found for path: %s", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	y.logger.Loglnf("[Yago] [%s] prepare handler for path: %s", handlerName, r.URL.Path)
	handler.ServeHTTP(w, r)
}

func (y *Yago) findHandler(p string) (http.Handler, string) {

	y.mu.RLock()
	if h, ok := y.paths[p]; ok {
		y.mu.RUnlock()
		return h.Handler(), h.Type()
	}
	y.mu.RUnlock()

	y.mu.Lock()
	defer y.mu.Unlock()

	for _, h := range y.handlers {
		if !strings.HasPrefix(p, h.Pattern()) {
			continue
		}
		y.paths[p] = h
		return h.Handler(), h.Type()
	}

	y.paths[p] = nil
	return nil, ""
}
