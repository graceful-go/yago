package yago

import (
	"context"
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

var (
	_ym    *YagoMessage = new(YagoMessage)
	_yc    *YagoContext = new(YagoContext)
	_ye    *YagoError   = new(YagoError)
	_codec YagoCodeC    = &YagoJsonCodec{}
)

type YagoAPIWrapper struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type YagoApiServerConfig struct {
	Route   string
	Timeout int
}

type YagoApiServer struct {
	c        *YagoApiServerConfig
	handlers map[string]*YagoApiHandler
	logger   Logger
}

func NewYagoApiServer(c *YagoApiServerConfig) (*YagoApiServer, error) {
	return &YagoApiServer{
		c:        c,
		handlers: make(map[string]*YagoApiHandler),
		logger:   &DefaultLogger{},
	}, nil
}

func (y *YagoApiServer) Type() string {
	return "YagoApiServer"
}

func (y *YagoApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	y.logger.Loglnf("[YagoApiServer] recv request: %s", r.URL.Path)

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

	handler, ok := y.handlers[strings.TrimPrefix(yc.serviceName, y.c.Route)]

	if !ok {
		y.logger.Loglnf("[YagoApiServer] Handle fail, handler not found for [%s]", yc.serviceName)
		yc.writeJson(&YagoAPIWrapper{
			Code: CodeYagoAPIServiceNotFound,
			Msg:  "service not found",
		})
		return
	}

	param, err := handler.packIn(yc.body)
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
	return y
}

func (y *YagoApiServer) Pattern() string {
	return y.c.Route
}

func (y *YagoApiServer) Register(serviceName string, handler interface{}) error {

	if _, ok := y.handlers[serviceName]; ok {
		return errors.New("duplicate service name registed:" + serviceName)
	}

	if hType := reflect.TypeOf(handler); hType.Kind() != reflect.Func {
		return errors.New("handler is not yago handler func")
	}

	h := &YagoApiHandler{
		fn: reflect.ValueOf(handler),
	}
	if err := h.init(); err != nil {
		return errors.New("invalid handler implementation for yago api handler")
	}
	y.handlers[serviceName] = h
	y.logger.Loglnf("[YagoApiServer] Register service succ for: %s/%s", y.Pattern(), serviceName)
	return nil
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
