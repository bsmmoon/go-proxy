package main

import (
	"os"
	"strconv"

	"github.com/bsmmoon/go-proxy/pkg/proxy"
	"github.com/bsmmoon/go-proxy/tool/logger"
)

var forever = make(chan int)

func main() {
	logger.SetCmd("Proxy")
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	options := proxy.ProxyOptions{
		SeleniumDriverPath: os.Getenv("SELENIUM_PATH"),
		GeckoDriverPath:    os.Getenv("GECKO_PATH"),
		Port:               port,
	}
	logger.INFO("\nPort: %v\nSelenium: %v\nGecko: %v", "8089", options.GeckoDriverPath, options.SeleniumDriverPath)

	go func() {
		proxy.Proxy(options)
	}()
	go func() {
		proxy.Selenium(options)
	}()

	<-forever
}
