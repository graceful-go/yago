package main

import (
	"fmt"
	"strings"
	"time"
)

type PageRouter struct {
	Routers []string
}

type Metadata struct {
	AppName    string
	AppVersion string
}

type SearchData struct {
}

type FooterClockData struct {
	Year    string
	Month   string
	Day     string
	DayTime string
	Ts      int64
}

type FooterData struct {
	Clock *FooterClockData
}

type PageData struct {
	PageRouter *PageRouter
	Metadata   *Metadata
	Search     *SearchData
	Footer     *FooterData
}

func NewClockData() *FooterClockData {
	n := time.Now()
	f := &FooterClockData{
		Year:    fmt.Sprintf("%02d", n.Year()),
		Month:   fmt.Sprintf("%02d", int(n.Month())),
		Day:     fmt.Sprintf("%02d", n.Day()),
		DayTime: n.Format("15:04:05"),
		Ts:      n.Unix(),
	}
	return f
}

func NewPageData(uPath string) *PageData {
	p := &PageData{}
	p.Metadata = meta
	p.PageRouter = &PageRouter{
		Routers: strings.Split(uPath, "/"),
	}
	p.Footer = &FooterData{
		Clock: NewClockData(),
	}
	return p
}
