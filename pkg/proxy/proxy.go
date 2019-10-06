package proxy

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bsmmoon/go-proxy/tool/logger"
	"github.com/elazarl/goproxy"
)

// Options Options
type Options struct {
	Port   int
	Filter Filter
}

// Filter Filter
type Filter struct {
	ContentType string
	Filename    string
}

func containAny(wantedTokens []string, target string) bool {
	for _, wantedToken := range wantedTokens {
		if strings.Contains(target, wantedToken) {
			return true
		}
	}
	return false
}

// Proxy proxy
func Proxy(options Options) {
	var err error
	logger.INFO("Starting Proxy: %v", options)
	proxy := goproxy.NewProxyHttpServer()

	sessionID := time.Now().Unix()
	filepath := fmt.Sprintf("./output/%v", sessionID)
	err = os.MkdirAll(filepath, os.ModePerm)
	if err != nil {
		logger.ERROR("Something went wrong creating filepath: %v", err.Error())
	}

	wantedTypes := strings.Split(options.Filter.ContentType, ",")

	proxy.OnRequest().DoFunc(
		func(response *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			return response, nil
		})

	proxy.OnResponse().DoFunc(
		func(response *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			if response == nil {
				logger.WARNING("No response given, ignored")
				return nil
			}

			var body []byte

			// TODO: use go-routine so that it returns response without waiting for copying the data
			defer func() {
				logger.INFO("Header: %v", response.Header)

				contentType := response.Header.Get("Content-Type")
				if !containAny(wantedTypes, contentType) {
					return
				}

				logger.INFO("RequestURI: %v", response.Request.URL.RequestURI())
				uriTokens := strings.Split(response.Request.URL.RequestURI(), "/")
				if len(uriTokens) == 0 {
					return
				}

				if len(body) == 0 {
					logger.INFO("Empty body, ignored")
					return
				}

				filename := uriTokens[len(uriTokens)-1]
				logger.INFO("Copying the body.. Filepath: %v, Filename: %v", filepath, filename)
				if len(filename) == 0 {
					logger.WARNING("Empty filename, ignored")
					return
				}
				err := ioutil.WriteFile(fmt.Sprintf("%v/%v", filepath, filename), body, 0644)
				if err != nil {
					logger.ERROR("Something went wrong writing to file: %v", err.Error())
				}
			}()

			body = copyBodyFromResponse(response)
			return response
		})

	err = http.ListenAndServe(fmt.Sprintf(":%v", options.Port), proxy)
	if err != nil {
		logger.FATAL(err.Error())
	}
}

func copyBodyFromResponse(response *http.Response) []byte {
	data, _ := ioutil.ReadAll(response.Body)                              // flush body into data
	response.Body = ioutil.NopCloser(bytes.NewBufferString(string(data))) // fill body with the copy of data
	return data
}
