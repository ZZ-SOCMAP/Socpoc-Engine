package proto

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/yanmengfei/spoce/library/utils"
)

const (
	DefaultUserAgent   = "Mozilla/5.0 (Windows NT 10.0; rv:78.0) Gecko/20100101 Firefox/78.0"
	DefaultContentType = "application/x-www-form-urlencoded"
)

var requestPool = sync.Pool{New: func() interface{} { return new(Request) }}

// SetupRequest create a request object
func SetupRequest(method, target, body string, headers map[string]string) (*Request, error) {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodDelete, http.MethodConnect, http.MethodOptions, http.MethodTrace:
	default:
		return nil, fmt.Errorf("invalid method by %s", method)
	}
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	var request = requestPool.Get().(*Request)
	request.Method = method
	request.Url = SetupURL(u)
	request.Headers = make(map[string]string, len(headers))
	for key := range headers {
		request.Headers[key] = headers[key]
	}
	if _, ok := request.Headers["User-Agent"]; !ok {
		request.Headers["User-Agent"] = DefaultUserAgent
	}
	if _, ok := request.Headers["Content-Type"]; !ok {
		request.Headers["Content-Type"] = DefaultContentType
		request.ContentType = DefaultContentType
	}
	if len(body) > 0 {
		if request.Body, err = hex.DecodeString(body); err != nil {
			request.Body = utils.StrToBytes(body)
		}
		request.Headers["Content-Length"] = strconv.Itoa(len(request.Body))
	}
	return request, nil
}

// ReleaseRequest recycle a request object
func ReleaseRequest(request *Request) {
	if request != nil {
		request.Reset()
		requestPool.Put(request)
	}
}

// ToFasthttp to fasthttp request struct
func (x *Request) ToFasthttp() *fasthttp.Request {
	var origin = fasthttp.AcquireRequest()
	origin.Header.SetMethod(x.Method)
	origin.SetRequestURI(x.Url.ToString())
	origin.SetBody(x.Body)
	for key, value := range x.Headers {
		origin.Header.Set(key, value)
	}
	return origin
}
