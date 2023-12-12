package main

import "time"

const (
	AppVersion string = "0.0.1"
	AppName    string = "办公助手"
)

var (
	meta = &Metadata{AppName: AppName, AppVersion: AppVersion + "." + time.Now().Format("200601021504")}
)
