package yago

type YagoMessage interface {
	String() string
}

type YagoError interface {
	Error() string
}
