package yago

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type YaogoTemplateHandler func(ctx *YagoContext) (data interface{}, err error)

type PageLayoutConfig struct {
	ServiceName string   `json:"name"`
	Templates   []string `json:"templates"`
}

// YagoTemplateConfig
// YagoTemplateConfig is used to bind template page for yago server
type YagoTemplateConfig struct {

	// Route is pattern prefix to serve for http server
	Route string `json:"route"`

	// LayoutDir is the local path that all template files defined
	LayoutDir string `json:"layoutDir"`

	// BaseLayouts will be rendered everytime
	BaseLayouts []string `json:"baseLayouts"`

	// PageLayouts
	PageLayouts []*PageLayoutConfig `json:"pageLayouts"`

	// Timeout
	Timeout int `json:"timeout"`
}

type YagoTemplateServer struct {
	hds       map[string]YaogoTemplateHandler
	renders   map[string]*YagoRender
	c         *YagoTemplateConfig
	logger    Logger
	bindFuncs map[string]interface{}
}

func NewYagoTemplateServer(c *YagoTemplateConfig) (*YagoTemplateServer, error) {
	yServer := &YagoTemplateServer{
		c:         c,
		bindFuncs: make(map[string]interface{}),
		logger:    &DefaultLogger{},
		hds:       make(map[string]YaogoTemplateHandler),
		renders:   make(map[string]*YagoRender),
	}

	return yServer, nil
}

// BindFuncs
func (y *YagoTemplateServer) BindFuncs(bindFuncs map[string]interface{}) {
	y.bindFuncs = bindFuncs
}

// ServeHTTP
func (y *YagoTemplateServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
	yc.Context = ctx

	fs := strings.Split(p, "/")
	if len(fs) >= 3 {
		yc.route = fs[1]
		yc.serviceName = fs[2]
	}

	y.logger.Loglnf("[YagoTemplateServer] Handle HTTP Request for [%s] %s, service: %s", method, p, yc.serviceName)

	y.Handle(yc)
}

func (y *YagoTemplateServer) Pattern() string {
	return y.c.Route
}

func (y *YagoTemplateServer) Handler() http.Handler {
	return y
}

func (y *YagoTemplateServer) Handle(ctx *YagoContext) {

	y.logger.Loglnf("[YagoTemplateServer] Handle request: %+v", ctx.query)

	hd, render, err := y.findHandler(ctx.serviceName)
	if err != nil {
		y.logger.Loglnf("[YagoTemplateServer] Handle HTTP Request fail for [%s] %s", ctx.serviceName, ctx.path)
		ctx.w.WriteHeader(http.StatusNotFound)
		return
	}

	if hd == nil {
		y.logger.Loglnf("[YagoTemplateServer] Handle HTTP Request fail for [%s] %s, empty handler", ctx.serviceName, ctx.path)
		ctx.w.WriteHeader(http.StatusNotFound)
		return
	}

	switch ctx.r.Method {
	case http.MethodGet:
		renderData, err := hd(ctx)
		if err != nil {
			y.logger.Log("[YagoTemplateServer] Handle HTTP GET Request fail for "+ctx.path, "logic handle fail")
			return
		}

		if err := render.Render(ctx, renderData); err != nil {
			y.logger.Log("[YagoTemplateServer] Handle HTTP Request fail for "+ctx.path, "render fail"+err.Error())
			return
		}
		return
	}

	// case http.MethodPost:
	// 	data, code := hd.Post(ctx)
	// 	if code != 0 {
	// 		y.logger.Loglnf("[YagoTemplateServer] Handle HTTP POST Request fail for %s, code: %d", ctx.path, code)
	// 		ctx.writeResponseStatus(code)
	// 		return
	// 	}
	// 	bs, _ := json.Marshal(data)
	// 	ctx.w.Write(bs)
	// 	return
	// }
	// y.logger.Log("[YagoTemplateServer] Handle HTTP Request fail for "+ctx.path, "handler not found")
	ctx.writeResponseStatus(http.StatusNotFound)
}

func (y *YagoTemplateServer) findHandler(serviceName string) (h YaogoTemplateHandler, r *YagoRender, e error) {

	hd, hOk := y.hds[serviceName]
	render, rOk := y.renders[serviceName]

	isNotFound := !hOk || !rOk

	if isNotFound {
		return nil, nil, errors.New("handler or render not found")
	}

	h = hd
	r = render
	e = nil
	return
}

func (y *YagoTemplateServer) register(serviceName string, handler YaogoTemplateHandler, render *YagoRender) error {
	if y.hds == nil {
		y.logger.Log("[YagoTemplateServer] Regist handler fail for " + serviceName)
		return nil
	}

	y.hds[serviceName] = handler
	y.renders[serviceName] = render
	y.logger.Log("[YagoTemplateServer] RegisterHandler succ for " + serviceName)
	return nil
}

func (y *YagoTemplateServer) parseQuery(query string) map[string]string {
	r := make(map[string]string, 0)
	values, err := url.ParseQuery(query)
	if err != nil {
		y.logger.Log("[YagoTemplateServer] ParseQuery fail for queryString: " + query)
		return r
	}
	for k, v := range values {
		r[k] = v[0]
	}
	return r
}

func (y *YagoTemplateServer) Type() string {
	return "YagoTemplateServer"
}

// Register
func (y *YagoTemplateServer) Register(service string, handler YaogoTemplateHandler) error {

	tmpls := y.getBindTemplates(service)

	render, err := NewRenderWithTemplates(tmpls, y.bindFuncs)
	if err != nil {
		y.logger.Log("[YagoServer] RegisterRouter fail with binding templates:", tmpls, "err is "+err.Error())
		return err
	}

	y.logger.Log("[YagoServer] RegisterRouter succ with binding templates:", tmpls)

	return y.register(service, handler, render)
}

func (y *YagoTemplateServer) getBindTemplates(serviceName string) []string {

	t := []string{}
	b := []string{}
	if y.c.BaseLayouts != nil && len(y.c.BaseLayouts) > 0 {
		b = append(b, y.c.BaseLayouts...)
	}
	if y.c.PageLayouts == nil {
		return y.getBindTemplatesWithDistDir(append(t, b...))
	}
	for _, v := range y.c.PageLayouts {
		if v.ServiceName == serviceName {
			t = append(t, v.Templates...)
			break
		}
	}
	return y.getBindTemplatesWithDistDir(append(t, b...))
}

func (y *YagoTemplateServer) getBindTemplatesWithDistDir(t []string) []string {
	if y.c.LayoutDir == "" {
		return t
	}
	r := []string{}
	for _, v := range t {
		r = append(r, y.c.LayoutDir+"/"+v)
	}
	return r
}
