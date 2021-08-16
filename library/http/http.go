package http

import (
	"crypto/tls"
	"time"

	"github.com/valyala/fasthttp"
)

// FasthttpClient interface
type FasthttpClient interface {
	// DoTimeout 不重定向
	DoTimeout(req *fasthttp.Request, resp *fasthttp.Response, t time.Duration) error
	// DoRedirects 跟随重定向
	DoRedirects(req *fasthttp.Request, resp *fasthttp.Response, maxRedirectsCount int) error
}

var (
	client      FasthttpClient // http client
	dialtimeout time.Duration  // set timeout
	maxredirect int            // max redirects
)

// Setup setup http client
func Setup(timeout time.Duration, redirects int) {
	dialtimeout = timeout * time.Second
	maxredirect = redirects
	client = &fasthttp.Client{
		MaxConnDuration:          0,
		MaxConnWaitTimeout:       dialtimeout,
		TLSConfig:                &tls.Config{InsecureSkipVerify: true}, // 跳过证书验证
		NoDefaultUserAgentHeader: true,
		DisablePathNormalizing:   true,
		ReadTimeout:              dialtimeout,
		WriteTimeout:             dialtimeout,
	}
}

// Do send http request
func Do(request *fasthttp.Request, redirect bool) (response *fasthttp.Response, err error) {
	response = fasthttp.AcquireResponse()
	defer func() {
		request.SetConnectionClose()
		fasthttp.ReleaseRequest(request)
	}()
	if redirect {
		err = client.DoRedirects(request, response, maxredirect)
	} else {
		err = client.DoTimeout(request, response, dialtimeout)
	}
	if err != nil {
		return response, err
	}
	return response, nil
}
