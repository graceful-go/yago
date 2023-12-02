package yago

type Handler interface {
	Handle(ctx *YagoContext) (data interface{}, err error)
}
