package main

import "github.com/graceful-go/yago"

type TodoHandler struct{}

func (d *TodoHandler) Handle(ctx *yago.YagoContext) (data interface{}, err error) {
	data = NewPageData(ctx.Path())
	return data, nil
}

type MainHandler struct{}

func (d *MainHandler) Handle(ctx *yago.YagoContext) (data interface{}, err error) {
	data = NewPageData(ctx.Path())
	return data, nil
}

type SettingHandler struct{}

func (d *SettingHandler) Handle(ctx *yago.YagoContext) (data interface{}, err error) {
	data = NewPageData(ctx.Path())
	return data, nil
}
