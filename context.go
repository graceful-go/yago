package yago

import (
	"context"
	"encoding/json"
	"net/http"
)

type YagoContext struct {
	path string

	// route is a part of req.Path
	// eg: http://host:port/a/apidemo
	// and route is /a, this is spectified by server option inition
	// http Path is equals: route + serviceName
	route string

	// serviceName is a part of req.Path
	// eg: http://host:port/a/apidemo
	// and serviceName is apidemo, this is spectified by server option inition
	// http Path is equals: route + serviceName
	serviceName string

	// query is http request query
	query map[string]string

	// body is http request body
	body []byte

	w http.ResponseWriter
	r *http.Request

	context.Context
}

func (y *YagoContext) writeResponseStatus(code int) {
	y.w.WriteHeader(code)
}

func (y *YagoContext) writeJson(data interface{}) {
	bs, _ := json.Marshal(data)
	y.w.Header().Set("Content-Type", "application/json")
	y.w.Write(bs)
}

func (y *YagoContext) ServiceName() string {
	return y.serviceName
}

func (y *YagoContext) Query(key string) string {
	if y.query == nil {
		return ""
	}
	return y.query[key]
}

func (y *YagoContext) Path() string {
	return y.path
}
