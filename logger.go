package yago

import "fmt"

type Logger interface {
	Log(msg ...interface{})
	Loglnf(format string, args ...interface{})
}

type DefaultLogger struct{}

func (d *DefaultLogger) Log(msgs ...interface{}) {
	fmt.Println(msgs...)
}

func (d *DefaultLogger) Loglnf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
