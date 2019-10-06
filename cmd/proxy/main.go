package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/bsmmoon/go-proxy/pkg/proxy"
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

	<-forever
}
