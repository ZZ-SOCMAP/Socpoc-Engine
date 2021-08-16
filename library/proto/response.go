package proto

import (
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
	"sync"
)

var responsePool = sync.Pool{New: func() interface{} { return new(Response) }}

// SetupResponse create a response object
func SetupResponse(response *fasthttp.Response, req *fasthttp.Request) (*Response, error) {
	u, err := url.Parse(req.URI().String())
	if err != nil {
		return nil, err
	}
	resp := responsePool.Get().(*Response)
	resp.Url = SetupURL(u)
	resp.Headers = make(map[string]string)
	resp.Status = int32(response.StatusCode())
	headerContent := response.Header.String()
	headers := strings.Split(headerContent, "\r\n")
	resp.ContentType = string(response.Header.Peek("Content-Type"))
	resp.Body = make([]byte, len(response.Body()))
	copy(resp.Body, response.Body())
	for i := 0; i < len(headers); i++ {
		values := strings.SplitN(headers[i], ":", 2)
		if len(values) != 2 {
			continue
		}
		resp.Headers[strings.ToLower(values[0])] = strings.TrimLeft(values[1], " ")
	}
	return resp, nil
}

// RecycleResponse recycle a response object
func RecycleResponse(response *Response) {
	if response != nil {
		response.Reset()
		responsePool.Put(response)
	}
}
