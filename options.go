package yago

type Option func(*Yago)

func WithConfigFile(path string) Option {
	return func(y *Yago) {
	}
}

func WithConfig(yc *YagoConfig) Option {
	return func(y *Yago) {
		y.yc = yc
	}
}

func WithFuncs(funcName string, bindFunc interface{}) Option {
	return func(y *Yago) {
		y.bindFuncs[funcName] = bindFunc
	}
}
