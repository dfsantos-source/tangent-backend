package main

import (
	"github.com/dfsantos-source/tangent-backend/http"
)

type Main struct {
	HTTPServer *http.Server
}

func CreateMain() *Main {
	return &Main{
		HTTPServer: http.CreateServer(),
	}
}

func (m *Main) Run() {
	m.HTTPServer.RunServer()
}

func main() {
	m := CreateMain()
	m.Run()
}
