package yago

type ServerConfig struct {
	Port    uint32 `json:"port" default:"8080"`
	Timeout uint32 `json:"timeout" default:"1000"`
}

type PageLayoutConfig struct {
	Path      string   `json:"path"`
	Templates []string `json:"templates"`
	Method    string   `json:"method"`
}

type BaseLayoutConfig struct {
	Templates []string `json:"templates"`
}

type YagoPageConfig struct {
	LayoutDir   string              `json:"baseDir"`
	AssetDir    string              `json:"assetDir"`
	BaseLayouts *BaseLayoutConfig   `json:"baseLayouts"`
	PageLayouts []*PageLayoutConfig `json:"pageLayouts"`
}

// YagoConfig
type YagoConfig struct {
	Server *ServerConfig   `json:"server"`
	Pages  *YagoPageConfig `json:"pages"`
}

func (y *YagoConfig) GetBindMethod(router string) string {
	if y.Pages.BaseLayouts == nil {
		return ""
	}
	for _, v := range y.Pages.PageLayouts {
		if v.Path == router {
			return v.Method
		}
	}
	return ""
}

func (y *YagoConfig) GetBindTemplates(router string) []string {
	t := []string{}
	b := []string{}
	if y.Pages.BaseLayouts != nil && len(y.Pages.BaseLayouts.Templates) > 0 {
		b = append(b, y.Pages.BaseLayouts.Templates...)
	}
	if y.Pages.PageLayouts == nil {
		return y.GetBindTemplatesWithDistDir(append(t, b...))
	}
	for _, v := range y.Pages.PageLayouts {
		if v.Path == router {
			t = append(t, v.Templates...)
			break
		}
	}
	return y.GetBindTemplatesWithDistDir(append(t, b...))
}

func (y *YagoConfig) GetBindTemplatesWithDistDir(t []string) []string {
	if y.Pages.LayoutDir == "" {
		return t
	}
	r := []string{}
	for _, v := range t {
		r = append(r, y.Pages.LayoutDir+"/"+v)
	}
	return r
}
