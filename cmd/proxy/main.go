package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/bsmmoon/go-proxy/pkg/proxy"
	"github.com/bsmmoon/go-proxy/pkg/seleniumwrapper"
	"github.com/bsmmoon/go-proxy/tool/logger"
)

var forever = make(chan int)

func main() {
	logger.SetCmd("Proxy")

	contentTypeFlag := flag.String("content-type", "", "ex. ")

	flag.Parse()
	contentType := *contentTypeFlag

	proxyPort, _ := strconv.Atoi(os.Getenv("PROXY_PORT"))
	go func() {
		proxy.Proxy(proxy.Options{
			Port: proxyPort,
			Filter: proxy.Filter{
				ContentType: contentType,
			},
		})
	}()

	seleniumPort, _ := strconv.Atoi(os.Getenv("SELENIUM_PORT"))
	go func() {
		seleniumwrapper.Selenium(seleniumwrapper.Option{
			SeleniumDriverPath: os.Getenv("SELENIUM_PATH"),
			GeckoDriverPath:    os.Getenv("GECKO_PATH"),
			ChromeDriverPath:   os.Getenv("CHROME_PATH"),
			Port:               seleniumPort,
			ProxyPort:          proxyPort,
			Browser:            seleniumwrapper.CHROME,
		})
	}()

	<-forever
}
