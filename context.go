package yago

import "net/http"

type YagoContext struct {
	path   string
	query  map[string]string
	body   []byte
	method string

	w http.ResponseWriter
	r *http.Request
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
