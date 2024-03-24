package yago

type Option func(*Yago) error

func WithConfig(yc *YagoConfig) Option {
	return func(y *Yago) error {
		y.cfg = yc
		return nil
	}
}

func WithFileServer(fsServer *YagoFileServer) Option {
	return func(y *Yago) error {
		y.handlers = append(y.handlers, fsServer)
		y.paths[fsServer.Pattern()] = fsServer
		return nil
	}
}

func WithTemplateServer(tServer *YagoTemplateServer) Option {
	return func(y *Yago) error {
		y.handlers = append(y.handlers, tServer)
		y.paths[tServer.Pattern()] = tServer
		return nil
	}
}

func WithApiServer(aServer *YagoApiServer) Option {
	return func(y *Yago) error {
		y.handlers = append(y.handlers, aServer)
		y.paths[aServer.Pattern()] = aServer
		return nil
	}
}
