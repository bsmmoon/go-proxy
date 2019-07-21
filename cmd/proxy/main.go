package main

import (
	"os"
	"strconv"

	"github.com/bsmmoon/go-proxy/pkg/proxy"
	"github.com/bsmmoon/go-proxy/tool/logger"
)

func main() {
	logger.SetCmd("Proxy")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	proxy.HelloWorld(proxy.ProxyOptions{
		SeleniumDriverPath: os.Getenv("SELENIUM_PATH"),
		GeckoDriverPath:    os.Getenv("GECKO_PATH"),
		Port:               port,
	})
}
