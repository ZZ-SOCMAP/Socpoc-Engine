package proto

import (
	"bytes"
	"fmt"
	"net/url"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/yanmengfei/poc-engine-xray/library/http"
	"github.com/yanmengfei/poc-engine-xray/library/utils"
)

var reverseDomain, reverseToken string

func SetupReverse(domain, token string) {
	reverseDomain = domain
	reverseToken = token
}

func NewReverse() *Reverse {
	if reverseDomain == "" {
		return &Reverse{}
	}
	var flag = utils.RandLetterNumbers(8)
	urlStr := fmt.Sprintf("http://%s.%s", flag, reverseDomain)
	u, _ := url.Parse(urlStr)
	return &Reverse{
		Url:                SetupURL(u),
		Flag:               flag,
		Domain:             u.Hostname(),
		Ip:                 "",
		IsDomainNameServer: false,
	}
}

// ceye.io api
const (
	dnsVerifyTemp  = "http://api.ceye.io/v1/records?token=%s&type=dns&filter=%s"
	httpVerifyTemp = "http://api.ceye.io/v1/records?token=%s&type=http&filter=%s"
)

// VerifyReverse 验证反连平台
func VerifyReverse(r *Reverse, timeout int64) bool {
	if reverseToken == "" {
		return false
	}
	time.Sleep(time.Second * time.Duration(timeout))                      // 延迟 x 秒获取结果
	if getReverseResp(fmt.Sprintf(dnsVerifyTemp, reverseToken, r.Flag)) { //check dns
		return true
	} else {
		if getReverseResp(fmt.Sprintf(httpVerifyTemp, reverseToken, r.Flag)) { //	check request
			return true
		}
	}
	return false
}

// getReverseResp 发送请求
func getReverseResp(verifyUrl string) bool {
	notExist := []byte(`"data": []`)
	var origin = fasthttp.AcquireRequest()
	origin.Header.SetMethod(fasthttp.MethodGet)
	origin.SetRequestURI(verifyUrl)
	response, err := http.Do(origin, false)
	defer func() {
		response.SetConnectionClose()
		fasthttp.ReleaseResponse(response)
	}()
	if err != nil {
		return false
	}
	if !bytes.Contains(response.Body(), notExist) {
		return true
	}
	return false
}
