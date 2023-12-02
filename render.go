package yago

import (
	"errors"
	"html/template"
	"path"
)

type YagoRender struct {
	t *template.Template
}

func NewRenderWithTemplates(tps []string, bindFuncs map[string]interface{}) (*YagoRender, error) {

	if len(tps) == 0 {
		return nil, errors.New("empty template found")
	}

	t, err := template.New(path.Base(tps[0])).Funcs(bindFuncs).ParseFiles(tps...)
	if err != nil {
		return nil, err
	}

	return &YagoRender{t: t}, nil
}

func NewRender(t *template.Template) *YagoRender {
	return &YagoRender{t: t}
}

func (y *YagoRender) Render(ctx *YagoContext, data interface{}) error {
	return y.t.Execute(ctx.w, data)
}
