package proxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/bsmmoon/go-proxy/tool/logger"
	"github.com/elazarl/goproxy"
	"github.com/tebeka/selenium"
)

type ProxyOptions struct {
	SeleniumDriverPath string
	GeckoDriverPath    string
	Port               int
}

// Proxy proxy
func Proxy(options ProxyOptions) {
	var err error
	logger.INFO("Starting Proxy")
	proxy := goproxy.NewProxyHttpServer()

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			r.Header.Set("X-GoProxy", "yxorPoG-X")
			return r, nil
		})

	proxy.OnResponse().DoFunc(
		func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			var body []byte
			defer func() {
				contentType := r.Header.Get("Content-Type")
				logger.INFO("Content-Type: %v", contentType)
				if contentType == "image/png" {
					logger.INFO("Header: %v", r.Header)
					logger.INFO("RequestURI: %v", r.Request.URL.RequestURI())
					uriTokens := strings.Split(r.Request.URL.RequestURI(), "/")
					if len(uriTokens) == 0 {
						return
					}
					filename := uriTokens[len(uriTokens)-1]
					logger.INFO("Filename: %v", filename)
					err := ioutil.WriteFile(fmt.Sprintf("./output/%v", filename), body, 0644)
					if err != nil {
						logger.WARNING("Something went wrong writing to file: %v", err.Error())
					}
				}
			}()

			body, _ = ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewBufferString(string(body)))
			return r
		})

	err = http.ListenAndServe(fmt.Sprintf(":%v", "8089"), proxy)
	if err != nil {
		logger.FATAL(err.Error())
	}
}

// Selenium selenium
func Selenium(options ProxyOptions) {
	logger.INFO("Starting Selenium")
	opts := []selenium.ServiceOption{
		// selenium.StartFrameBuffer(),        // Start an X frame buffer for the browser to run in. // xvfb not supported in MacOS?
		selenium.GeckoDriver(options.GeckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
		selenium.Output(os.Stderr),                    // Output debug information to STDERR.
	}
	service, err := selenium.NewSeleniumService(options.SeleniumDriverPath, options.Port, opts...)
	if err != nil {
		logger.FATAL(err.Error()) // panic is used only as an example and is not otherwise recommended.
	}
	defer service.Stop()

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{
		"browserName": "firefox",
		"proxy": selenium.Proxy{
			Type: selenium.Manual,
			HTTP: "localhost:8089",
		},
	}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", options.Port))
	if err != nil {
		logger.FATAL(err.Error())
	}
	defer wd.Quit()

	// Navigate to the simple playground interface.
	if err := wd.Get("http://eu.httpbin.org/"); err != nil {
		panic(err)
	}

	logger.INFO("insert y value here: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	logger.INFO(input.Text())
}
