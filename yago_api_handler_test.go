package yago

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DemoReq struct {
	Field string
}

func (d *DemoReq) String() string {
	return ""
}

type DemoRsp struct {
	Field string
}

func (d *DemoRsp) String() string {
	return ""
}

func TestYagoHandlerPackIn(t *testing.T) {
	handler := &YagoApiHandler{in: reflect.TypeOf(&DemoReq{})}
	rsp, err := handler.packIn([]byte(`{"Field":"hello"}`))
	assert.Equal(t, nil, err)
	assert.EqualValues(t, YagoMessage(&DemoRsp{Field: "hello"}), rsp)
}

func TestYagoHandlerInvoke(t *testing.T) {
	var fn = func(ctx *YagoContext, in *DemoReq) (*DemoRsp, error) {
		if in == nil {
			return nil, errors.New("nil in")
		}
		if in.Field == "Error" {
			return nil, errors.New("error")
		}
		if in.Field == "Nil" {
			return nil, nil
		}
		if in.Field == "Both" {
			return &DemoRsp{Field: "Both"}, errors.New("both error")
		}
		return &DemoRsp{Field: in.Field}, nil
	}

	var uts = []struct {
		Func         interface{}
		Param        YagoMessage
		ExpectErrNil bool
		ExpectRsp    YagoMessage
		Id           int
	}{
		{
			Func:         fn,
			Param:        &DemoReq{Field: "HelloWorld"},
			ExpectErrNil: true,
			ExpectRsp:    &DemoRsp{Field: "HelloWorld"},
			Id:           1,
		},
		{
			Func:         fn,
			Param:        nil,
			ExpectErrNil: false,
			ExpectRsp:    nil,
			Id:           2,
		},
		{
			Func:         fn,
			Param:        &DemoReq{Field: "Error"},
			ExpectErrNil: false,
			ExpectRsp:    nil,
			Id:           3,
		},
		{
			Func:         fn,
			Param:        &DemoReq{Field: "Nil"},
			ExpectErrNil: false,
			ExpectRsp:    nil,
			Id:           4,
		},
		{
			Func:         fn,
			Param:        &DemoReq{Field: "Both"},
			ExpectErrNil: false,
			ExpectRsp:    nil,
			Id:           5,
		},
	}
	for _, uc := range uts {
		yHandler := &YagoApiHandler{fn: reflect.ValueOf(uc.Func)}
		rsp, err := yHandler.invoke(&YagoContext{}, uc.Param)
		assert.Equal(t, uc.ExpectErrNil, err == nil)
		assert.EqualValues(t, uc.ExpectRsp, rsp)
	}
}

func TestYagoHandlerInit(t *testing.T) {

	var uts = []struct {
		Func           interface{}
		ExpectInitSucc bool
	}{
		{
			Func:           func(ctx *YagoContext, in *DemoReq) (*DemoRsp, error) { return nil, nil },
			ExpectInitSucc: true,
		},
		{
			Func:           func(ctx context.Context, in *DemoReq) (*DemoRsp, error) { return nil, nil },
			ExpectInitSucc: false,
		},
		{
			Func:           func(ctx *YagoContext, in string) (*DemoRsp, error) { return nil, nil },
			ExpectInitSucc: false,
		},
		{
			Func:           func(ctx *YagoContext, in *DemoReq) (string, error) { return "nil", nil },
			ExpectInitSucc: false,
		},
		{
			Func:           func(ctx *YagoContext, in *DemoReq) (*DemoRsp, string) { return nil, "" },
			ExpectInitSucc: false,
		},
		{
			Func:           func(ctx *YagoContext, in *DemoReq, in1 string) (*DemoRsp, string) { return nil, "" },
			ExpectInitSucc: false,
		},
		{
			Func:           func(ctx *YagoContext, in *DemoReq) *DemoRsp { return nil },
			ExpectInitSucc: false,
		},
	}

	for _, uc := range uts {
		yHandler := &YagoApiHandler{fn: reflect.ValueOf(uc.Func)}
		assert.Equal(t, uc.ExpectInitSucc, yHandler.init() == nil)
	}
}
