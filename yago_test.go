package yago

import (
	"testing"
)

type DemoReq struct {
	Field string
}

type DemoRsp struct {
	Field string
}

func (d *DemoReq) String() string {
	return ""
}

func (d *DemoRsp) String() string {
	return ""
}

func TestYagoHandler(t *testing.T) {

}
