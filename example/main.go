package main

import (
	"context"

	"github.com/graceful-go/yago"
)

func main() {
	yc := &yago.YagoConfig{
		Server: &yago.ServerConfig{
			Port:    8080,
			Timeout: 1000,
		},
		Pages: &yago.YagoPageConfig{
			LayoutDir: "static/layout",
			AssetDir:  "static/assets",
			BaseLayouts: &yago.BaseLayoutConfig{
				Templates: []string{
					"header.layout",
					"body.layout",
					"footer.layout",
				},
			},
			PageLayouts: []*yago.PageLayoutConfig{
				{
					Path:      "/",
					Templates: []string{"main.layout"},
					Method:    "*",
				},
			},
		},
	}

	ycs, err := yago.New(yago.WithConfig(yc))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	if err := ycs.RegisterRouter(ctx, "/", &DemoHandler{}); err != nil {
		panic(err)
	}

	ycs.Start(ctx)
}

type DemoHandler struct{}

func (d *DemoHandler) Handle(ctx *yago.YagoContext) (data interface{}, err error) {
	return nil, nil
}
