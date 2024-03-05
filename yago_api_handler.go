package yago

import (
	"errors"
	"reflect"
)

type YagoApiHandler struct {
	fn  reflect.Value
	in  reflect.Type
	out reflect.Type
}

func (y *YagoApiHandler) packIn(bs []byte) (YagoMessage, error) {

	param := reflect.New(y.in.Elem()).Interface()
	if err := _codec.Unmarshal(bs, param); err != nil {
		return nil, err
	}
	if p, ok := param.(YagoMessage); ok {
		return p, nil
	}
	return nil, errors.New("marshal fail, not YagoMessage found")
}

func (y *YagoApiHandler) invoke(yc *YagoContext, in YagoMessage) (YagoMessage, error) {

	if yc == nil || in == nil {
		return nil, errors.New("[YagoApiHandler] invoke fail, unexpected error occour, invoke params can not be nil value")
	}

	outs := y.fn.Call([]reflect.Value{
		reflect.ValueOf(yc),
		reflect.ValueOf(in),
	})
	if len(outs) != 2 {
		return nil, errors.New("[YagoApiHandler] invoke fail, unexpected error occour, response format error")
	}

	if err, ok := outs[1].Interface().(error); ok && err != nil {
		return nil, err
	}

	if outs[0].IsNil() && outs[1].IsNil() {
		return nil, errors.New("[YagoApiHandler] invoke fail, unexpected error occour, response format error")
	}

	rsp, ok := outs[0].Interface().(YagoMessage)
	if !ok {
		return nil, errors.New("[YagoApiHandler] invoke fail, unexpected error occour, response format error")
	}

	return rsp, nil
}

func (y *YagoApiHandler) init() error {

	if numIn := y.fn.Type().NumIn(); numIn != 2 {
		return errors.New("[YagoApiHandler] init fail(0x1), invalid handler implementation for yago api handler, params-in invalid")
	}

	if y.fn.Type().In(0) != (reflect.TypeOf(_yc)) {
		return errors.New("[YagoApiHandler] init fail(0x2), invalid handler implementation for yago api handler, params-in invalid")
	}

	if !y.fn.Type().In(1).ConvertibleTo(reflect.TypeOf(_ym).Elem()) {
		return errors.New("[YagoApiHandler] init fail(0x3), invalid handler implementation for yago api handler, params-in invalid")
	}

	if numOut := y.fn.Type().NumOut(); numOut != 2 {
		return errors.New("[YagoApiHandler] init fail(0x4), invalid handler implementation for yago api handler, params-out invalid")
	}

	if !y.fn.Type().Out(0).ConvertibleTo(reflect.TypeOf(_ym).Elem()) {
		return errors.New("[YagoApiHandler] init fail(0x5), invalid handler implementation for yago api handler, params-out invalid")
	}

	if !y.fn.Type().Out(1).ConvertibleTo(reflect.TypeOf(_ye).Elem()) {
		return errors.New("[YagoApiHandler] init fail(0x6), invalid handler implementation for yago api handler, params-out invalid, 6")
	}

	y.in = y.fn.Type().In(1)
	y.out = y.fn.Type().Out(0)

	return nil
}
