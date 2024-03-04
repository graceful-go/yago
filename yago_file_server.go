package yago

import (
	"errors"
	"net/http"
)

// YagoFileServerConfig
// serve static file server support for yago server
type YagoFileServerConfig struct {
	Dir   string `json:"dir"`
	Route string `json:"route"`
}

type YagoFileServer struct {
	fsConfig  *YagoFileServerConfig
	fsPath    string
	fsHandler http.Handler
}

func NewYagoFileServer(fsConfig *YagoFileServerConfig) (*YagoFileServer, error) {
	if fsConfig == nil || fsConfig.Dir == "" || fsConfig.Route == "" {
		return nil, errors.New("empty fsConfig is not allowed")
	}
	fsPath := "/" + fsConfig.Route + "/"
	fsHandler := http.StripPrefix(fsPath, http.FileServer(http.Dir("./"+fsConfig.Dir)))
	return &YagoFileServer{
		fsConfig:  fsConfig,
		fsHandler: fsHandler,
		fsPath:    fsPath,
	}, nil
}

func (y *YagoFileServer) Handler() http.Handler {
	return y.fsHandler
}

func (y *YagoFileServer) Pattern() string {
	return y.fsPath
}

func (y *YagoFileServer) Type() string {
	return "YagoFileServer"
}
