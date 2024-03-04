package yago

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

const (
	CodeYagoAPISucc            int = 0
	CodeYagoAPIReqReadError    int = -100
	CodeYagoAPIReqParseError   int = -101
	CodeYagoAPIInternalError   int = -102
	CodeYagoAPIServiceNotFound int = -110
)

type YagoAPIWrapper struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type YagoApiHandler struct {
	handler YagoApiHandlerFunc
	in      reflect.Type
}

type YagoApiHandlerFunc func(ctx context.Context, req interface{}) (interface{}, error)

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
	param := reflect.New(handler.in.Elem()).Interface()
	if err := json.Unmarshal(yc.body, param); err != nil {
		yc.writeJson(&YagoAPIWrapper{
			Code: CodeYagoAPIReqParseError,
			Msg:  "req param type not match",
		})
		return
	}
	rsp, err := handler.handler(yc, param)
	if err != nil {
		yc.writeJson(&YagoAPIWrapper{
			Code: CodeYagoAPIInternalError,
			Msg:  "handle error",
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

func (y *YagoApiServer) Register(serviceName string, handler YagoApiHandlerFunc) {

	if _, ok := y.handlers[serviceName]; ok {
		panic("duplicate service name registed:" + serviceName)
	}

	meta := reflect.ValueOf(handler)
	if numIn := meta.Type().NumIn(); numIn != 1 {
		panic("invalid handler implement for yago api handler")
	}
	paramType := meta.Type().In(0)

	h := &YagoApiHandler{
		handler: handler,
		in:      paramType,
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
