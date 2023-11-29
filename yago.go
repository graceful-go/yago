package yago

import (
	"context"
	"fmt"
	"net/http"
)

type Yago struct {
	yc     *YagoConfig
	ys     *YagoServer
	logger Logger
}

func New(opts ...Option) (*Yago, error) {

	logger := &DefaultLogger{}

	y := &Yago{logger: logger}
	for _, opt := range opts {
		opt(y)
	}

	ys := &YagoServer{}
	ys.hds = make(map[string]Handler)
	ys.renders = make(map[string]*YagoRender)
	ys.c = y.yc.Server
	ys.logger = logger
	y.ys = ys

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
	if y.yc.Pages.AssetDir != "" {
		assetPath := "/" + y.yc.Pages.AssetDir + "/"
		fsHandler := http.StripPrefix(assetPath, http.FileServer(http.Dir("./"+y.yc.Pages.AssetDir)))
		http.Handle(assetPath, fsHandler)
		y.logger.Log("[YagoServer] Server Start fs handler for: ", assetPath)
	}
	http.Handle("/", y.ys)
	y.logger.Log("[YagoServer] Server Start page handler for: /")
	y.logger.Log("[YagoServer] Server Startup at", y.yc.Server.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", y.ys.c.Port), nil)
}

// RegisterRouter
func (y *Yago) RegisterRouter(ctx context.Context, router string, handler Handler) error {

	tmpls := y.yc.GetBindTemplates(router)
	method := y.yc.GetBindMethod(router)

	render, err := NewRenderWithTemplates(tmpls)
	if err != nil {
		y.logger.Log("[YagoServer] RegisterRouter fail with binding templates:", tmpls, "err is "+err.Error())
		return err
	}

	y.logger.Log("[YagoServer] RegisterRouter succ with binding templates:", tmpls)

	return y.ys.Register(router, method, handler, render)
}
