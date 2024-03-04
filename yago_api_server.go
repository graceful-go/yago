package yago

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const (
	CodeYagoAPISucc            int = 0
	CodeYagoAPIReqReadError    int = -100004
	CodeYagoAPIReqParseError   int = -100003
	CodeYagoAPIInternalError   int = -100002
	CodeYagoAPIServiceNotFound int = -100001
)

type YagoAPIWrapper struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type YagoApiHandler struct {
	fn  reflect.Value
	in  reflect.Type
	out reflect.Type
}

func (y *YagoApiHandler) parse(bs []byte) (YagoMessage, error) {

	param := reflect.New(y.in.Elem()).Interface()
	if err := json.Unmarshal(bs, param); err != nil {
		return nil, err
	}
	if p, ok := param.(YagoMessage); ok {
		return p, nil
	}
	return nil, errors.New("marshal fail, not YagoMessage found")
}

func (y *YagoApiHandler) invoke(yc *YagoContext, in YagoMessage) (YagoMessage, error) {
	outs := y.fn.Call([]reflect.Value{
		reflect.ValueOf(yc),
		reflect.ValueOf(in),
	})
	if len(outs) != 2 {
		return nil, errors.New("unexpected error occour, response format error")
	}

	if err, ok := outs[1].Interface().(error); ok && err != nil {
		return nil, err
	}

	if rsp, ok := outs[0].Interface().(YagoMessage); ok {
		return rsp, nil
	}

	return nil, errors.New("unexpected error occour, response format error")
}

func (y *YagoApiHandler) init() error {

	if numIn := y.fn.Type().NumIn(); numIn != 2 {
		return errors.New("invalid handler implementation for yago api handler")
	}

	if !y.fn.Type().In(0).ConvertibleTo(reflect.TypeOf(&YagoContext{})) {
		return errors.New("invalid handler implementation for yago api handler")
	}

	if !y.fn.Type().In(1).Implements(reflect.TypeOf(new(YagoMessage))) {
		return errors.New("invalid handler implementation for yago api handler")
	}

	if numOut := y.fn.Type().NumOut(); numOut != 2 {
		return errors.New("invalid handler implementation for yago api handler")
	}

	if !y.fn.Type().Out(0).Implements(reflect.TypeOf(new(YagoMessage))) {
		return errors.New("invalid handler implementation for yago api handler")
	}

	if !y.fn.Type().Out(1).Implements(reflect.TypeOf(new(error))) {
		return errors.New("invalid handler implementation for yago api handler")
	}

	y.in = y.fn.Type().In(1)
	y.out = y.fn.Type().Out(0)

	return nil
}

type YagoApiServerConfig struct {
	Route   string
	Timeout int
}

type YagoApiServer struct {
	c        *YagoApiServerConfig
	handlers map[string]*YagoApiHandler
	mux      *http.ServeMux
	logger   Logger
}

func NewYagoApiServer(c *YagoApiServerConfig) (*YagoApiServer, error) {
	return &YagoApiServer{
		c:        c,
		handlers: make(map[string]*YagoApiHandler),
		mux:      http.NewServeMux(),
	}, nil
}

func (y *YagoApiServer) Type() string {
	return "YagoApiServer"
}

func (y *YagoApiServer) Handle(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(y.c.Timeout))
	defer cancel()

	queryString, _ := url.QueryUnescape(r.URL.RawQuery)
	queryParams := y.parseQuery(queryString)

	p := r.URL.Path
	method := r.Method

	yc := &YagoContext{}
	yc.path = p
	yc.query = queryParams
	yc.w = w
	yc.r = r
	yc.route = y.c.Route
	yc.Context = ctx
	yc.serviceName = strings.TrimPrefix(r.URL.Path, y.c.Route)

	y.logger.Loglnf("[YagoApiServer] Handle HTTP Request for [%s] %s", method, p)

	if r.Body != nil {
		bs, err := io.ReadAll(r.Body)
		if err != nil {
			y.logger.Loglnf("[YagoApiServer] Handle HTTP Request fail for [%s], err: %s", method, err.Error())
			yc.writeJson(&YagoAPIWrapper{
				Code: CodeYagoAPIReqReadError,
				Msg:  "read request body fail",
			})
			return
		}
		yc.body = bs
		if err := r.Body.Close(); err != nil {
			y.logger.Loglnf("[YagoApiServer] Handle HTTP Request fail for [%s], err: %s", method, err.Error())
			yc.writeJson(&YagoAPIWrapper{
				Code: CodeYagoAPIReqParseError,
				Msg:  "parse request body fail",
			})
			return
		}
	}

	y.invoke(yc)
}

func (y *YagoApiServer) invoke(yc *YagoContext) {

	handler, ok := y.handlers[yc.serviceName]
	if !ok {
		y.logger.Loglnf("[YagoApiServer] Handle fail, handler not found for [%s], err: %s", yc.serviceName)
		yc.writeJson(&YagoAPIWrapper{
			Code: CodeYagoAPIServiceNotFound,
			Msg:  "service not found",
		})
		return
	}

	param, err := handler.parse(yc.body)
	if err != nil {
		yc.writeJson(&YagoAPIWrapper{
			Code: CodeYagoAPIReqParseError,
			Msg:  "req param type not match",
		})
		return
	}
	rsp, err := handler.invoke(yc, param)
	if err != nil {
		yc.writeJson(&YagoAPIWrapper{
			Code: CodeYagoAPIInternalError,
			Msg:  "invoke error",
		})
		return
	}
	yc.writeJson(&YagoAPIWrapper{Data: rsp})
}

func (y *YagoApiServer) Handler() http.Handler {
	return y.mux
}

func (y *YagoApiServer) Pattern() string {
	return y.c.Route
}

func (y *YagoApiServer) Register(serviceName string, handler interface{}) {

	if _, ok := y.handlers[serviceName]; ok {
		panic("duplicate service name registed:" + serviceName)
	}

	if hType := reflect.TypeOf(handler); hType.Kind() != reflect.Func {
		panic("handler is not yago handler func")
	}

	h := &YagoApiHandler{
		fn: reflect.ValueOf(handler),
	}
	if err := h.init(); err != nil {
		panic("invalid handler implementation for yago api handler")
	}
	y.handlers[serviceName] = h
}

func (y *YagoApiServer) parseQuery(query string) map[string]string {
	r := make(map[string]string, 0)
	values, err := url.ParseQuery(query)
	if err != nil {
		y.logger.Log("[YagoApiServer] ParseQuery fail for queryString: " + query)
		return r
	}
	for k, v := range values {
		r[k] = v[0]
	}
	return r
}
