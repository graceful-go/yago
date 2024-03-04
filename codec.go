package yago

type YagoCodeC interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte, interface{}) error
	Name() string
}
