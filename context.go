package yago

import (
	"context"
	"encoding/json"
	"net/http"
)

type YagoContext struct {
	path   string
	query  map[string]string
	body   []byte
	method string

	w http.ResponseWriter
	r *http.Request

	context.Context
}

func (y *YagoContext) writeResponseStatus(code int) {
	y.w.WriteHeader(code)
}

func (y *YagoContext) Method() string {
	return y.method
}

func (y *YagoContext) ParseBody(dst interface{}) error {
	return json.Unmarshal(y.body, dst)
}

func (y *YagoContext) Body() []byte {
	return y.body
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
