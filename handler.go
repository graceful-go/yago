package yago

type Handler interface {
	Get(ctx *YagoContext) (data interface{}, err error)
	Post(ctx *YagoContext) (data interface{}, code int)
}
