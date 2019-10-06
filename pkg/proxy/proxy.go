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

func containAny(strs, del, target string) bool {
	for _, str := range strings.Split(strs, del) {
		if strings.Contains(target, str) {
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
				if containAny(options.Filter.ContentType, ",", contentType) {
					logger.INFO("Header: %v", r.Header)
					logger.INFO("RequestURI: %v", r.Request.URL.RequestURI())
					uriTokens := strings.Split(r.Request.URL.RequestURI(), "/")
					if len(uriTokens) == 0 {
						return
					}
					filepath := fmt.Sprintf("./output/%v", sessionID)
					filename := uriTokens[len(uriTokens)-1]
					logger.INFO("Filepath: %v, Filename: %v", filepath, filename)
					err = os.MkdirAll(filepath, os.ModePerm)
					if err != nil {
						logger.WARNING("Something went wrong creating filepath: %v", err.Error())
					}
					err = ioutil.WriteFile(fmt.Sprintf("%v/%v", filepath, filename), body, 0644)
					if err != nil {
						logger.WARNING("Something went wrong writing to file: %v", err.Error())
					}
				}
			}()

			body, _ = ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewBufferString(string(body)))
			return r
		})

	err = http.ListenAndServe(fmt.Sprintf(":%v", options.Port), proxy)
	if err != nil {
		logger.FATAL(err.Error())
	}
}
