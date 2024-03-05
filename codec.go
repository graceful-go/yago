package yago

import "encoding/json"

type YagoCodeC interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte, interface{}) error
	Name() string
}

type YagoJsonCodec struct{}

func (y *YagoJsonCodec) Marshal() ([]byte, error) {
	return json.Marshal(y)
}

func (y *YagoJsonCodec) Unmarshal(bs []byte, dst interface{}) error {
	return json.Unmarshal(bs, dst)
}

func (y *YagoJsonCodec) Name() string {
	return "json"
}
