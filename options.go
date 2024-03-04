package yago

type Option func(*Yago) error

func WithConfig(yc *YagoConfig) Option {
	return func(y *Yago) error {
		y.cfg = yc
		return nil
	}
}

func WithFileServer(fsConfig *YagoFileServerConfig) Option {
	return func(y *Yago) error {
		fsServer, err := NewYagoFileServer(fsConfig)
		if err != nil {
			y.logger.Log("[YagoServer] Server start fail due to:", err.Error())
			return err
		}
		y.handlers = append(y.handlers, fsServer)
		return nil
	}
}

func WithTemplateServer(tlConfig *YagoTemplateConfig) Option {
	return func(y *Yago) error {
		fsServer, err := NewYagoTemplateServer(tlConfig)
		if err != nil {
			y.logger.Log("[YagoServer] Server start fail due to:", err.Error())
			return err
		}
		y.handlers = append(y.handlers, fsServer)
		return nil
	}
}
