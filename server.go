package yago

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type YagoServer struct {
	hds     map[string]Handler
	renders map[string]*YagoRender

	c      *ServerConfig
	logger Logger
}

func (y *YagoServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(y.c.Timeout))
	defer cancel()

	queryString, _ := url.QueryUnescape(r.URL.RawQuery)
	queryParams := y.ParseQuery(ctx, queryString)

	p := r.URL.Path
	method := r.Method

	y.logger.Loglnf("[YagoServer] Handle HTTP Request for [%s] %s", method, p)

	yc := &YagoContext{}
	yc.path = p
	yc.query = queryParams
	yc.w = w
	yc.r = r
	yc.body = y.ParseBody(r)
	yc.method = method

	y.Handle(yc)
}

func (y *YagoServer) Handle(ctx *YagoContext) {

	hd, render, err := y.FindHandler(ctx.path, ctx.method)
	if err != nil {
		y.logger.Loglnf("[YagoServer] Handle HTTP Request fail for [%s] %s", ctx.method, ctx.path)
		ctx.w.WriteHeader(http.StatusNotFound)
		return
	}

	switch ctx.Method() {
	case http.MethodGet:
		renderData, err := hd.Get(ctx)
		if err != nil {
			y.logger.Log("[YagoServer] Handle HTTP GET Request fail for "+ctx.path, "logic handle fail")
			return
		}

		if err := render.Render(ctx, renderData); err != nil {
			y.logger.Log("[YagoServer] Handle HTTP Request fail for "+ctx.path, "render fail"+err.Error())
			return
		}
		return
	case http.MethodPost:
		data, code := hd.Post(ctx)
		if code != 0 {
			y.logger.Loglnf("[YagoServer] Handle HTTP POST Request fail for %s, code: %d", ctx.path, code)
			ctx.writeResponseStatus(code)
			return
		}
		bs, _ := json.Marshal(data)
		ctx.w.Write(bs)
		return
	}
	y.logger.Log("[YagoServer] Handle HTTP Request fail for "+ctx.path, "handler not found")
	ctx.writeResponseStatus(http.StatusNotFound)
}

func (y *YagoServer) FindHandler(path, method string) (h Handler, r *YagoRender, e error) {

	regKey := y.RegKey(path, method)
	hd, hOk := y.hds[regKey]
	render, rOk := y.renders[regKey]

	isNotFound := !hOk || !rOk

	if isNotFound && method != "*" {
		return y.FindHandler(path, "*")
	}

	if isNotFound {
		return nil, nil, errors.New("handler or render not found")
	}

	h = hd
	r = render
	e = nil
	return
}

func (y *YagoServer) RegKey(path, method string) string {
	return fmt.Sprintf("[%s] %s", method, path)
}

func (y *YagoServer) Register(path string, method string, handler Handler, render *YagoRender) error {
	if y.hds == nil {
		y.logger.Log("[YagoServer] Regist handler fail for " + path)
		return nil
	}

	regKey := y.RegKey(path, method)

	y.hds[regKey] = handler
	y.renders[regKey] = render
	y.logger.Log("[YagoServer] RegisterHandler succ for " + regKey)
	return nil
}

func (y *YagoServer) ParseBody(r *http.Request) []byte {
	if r.Body == nil {
		return []byte{}
	}
	bs, _ := io.ReadAll(r.Body)
	return bs
}

func (y *YagoServer) ParseQuery(ctx context.Context, query string) map[string]string {
	r := make(map[string]string, 0)
	values, err := url.ParseQuery(query)
	if err != nil {
		y.logger.Log("[YagoServer] ParseQuery fail for queryString: " + query)
		return r
	}
	for k, v := range values {
		r[k] = v[0]
	}
	return r
}
