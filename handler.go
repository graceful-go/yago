package yago

type Handler interface {
	Handle(ctx *YagoContext) (renderData interface{}, err error)
}
