package main

import (
	"github.com/bsmmoon/go-proxy/pkg/proxy"
	"github.com/bsmmoon/go-proxy/tool/logger"
)

func main() {
	logger.SetCmd("Proxy")
	proxy.HelloWorld()
}
