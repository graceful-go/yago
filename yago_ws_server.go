package yago

import "net/http"

type YagoWsServerConfig struct {
	Route string
}

type YagoWsServer struct {
	mux *http.ServeMux
	c   *YagoWsServerConfig
}

func NewYagoWsServer(c *YagoWsServerConfig) (*YagoWsServer, error) {

	y := &YagoWsServer{
		c:   c,
		mux: http.NewServeMux(),
	}

	y.mux.HandleFunc(c.Route, y.Handle)
	return y, nil
}

func (y *YagoWsServer) Handle(http.ResponseWriter, *http.Request) {

}

func (y *YagoWsServer) Type() string {
	return "YagoWsServer"
}

func (y *YagoWsServer) Handler() http.Handler {
	return y.mux
}

func (y *YagoWsServer) Pattern() string {
	return y.c.Route
}
