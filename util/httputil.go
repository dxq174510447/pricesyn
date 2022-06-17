package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"
)

type httpUtil struct {
}

func (h *httpUtil) IsTimeOutError(err error) bool {
	if err == nil {
		return false
	}
	if netError, err := err.(net.Error); err {
		// 判断是否超时
		if netError.Timeout() {
			return true
		}
	}
	return false
}

func (h *httpUtil) PostBody(ctx context.Context, client *http.Client, url string,
	requestBody interface{}, responseBody interface{}, header map[string]string) error {

	var reqBody []byte
	var err error

	if reflect.TypeOf(requestBody).Kind() == reflect.String {
		str := requestBody.(string)
		reqBody = []byte(str)
	} else {
		reqBody, err = json.Marshal(requestBody)
	}

	if err != nil {
		return err
	}
	reqReader := bytes.NewReader(reqBody)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, reqReader)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if len(header) > 0 {
		for k, v := range header {
			httpReq.Header.Set(k, v)
		}
	}

	response, err := client.Do(httpReq)

	if err != nil {
		return err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(res, responseBody)
	if err != nil {
		return err
	}
	return nil
}

func (h *httpUtil) Get(ctx context.Context, client1 *http.Client, url string, header map[string]string) ([]byte, error) {

	client := client1
	if client == nil {
		client = h.GetHttpClient(10, 30, "")
	}

	var err error
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if len(header) > 0 {
		for k, v := range header {
			httpReq.Header.Set(k, v)
		}
	}

	response, err := client.Do(httpReq)

	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// 创建默认httpclient 插入第三方api中 实现请求日志打印，（监控等）
func (h *httpUtil) GetHttpClient(timeout int, keepalive int, proxyUrl string) *http.Client {
	roundTripper := h.GetHttpRoundTripper(timeout, keepalive, proxyUrl)
	var client = http.Client{
		Transport: http.RoundTripper(roundTripper),
	}
	return &client
}

func (h *httpUtil) GetHttpRoundTripper(timeout int, keepalive int, proxyUrl string) http.RoundTripper {
	var transport *http.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
			KeepAlive: time.Duration(keepalive) * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   20 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if proxyUrl != "" {
		pURL, _ := url.Parse(proxyUrl)
		transport.Proxy = http.ProxyURL(pURL)
	}

	tripper := &platformRoundTripper{
		Transport:  *transport,
		timeoutKey: make(map[string]int),
	}

	return http.RoundTripper(tripper)
}

type platformTimeoutError struct {
}

func (p platformTimeoutError) Error() string {
	return "mock timeout error"
}

func (p platformTimeoutError) Timeout() bool {
	return true
}

func (p platformTimeoutError) Temporary() bool {
	return false
}

type platformRoundTripper struct {
	http.Transport
	timeoutKey map[string]int
}

func (t *platformRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// http头
	m := req.Context().Value("HTTP_HEAD")
	if m != nil {
		var hs map[string]string = m.(map[string]string)
		for k, v := range hs {
			req.Header.Set(k, v)
		}
	}

	// 是否打印原始报文
	var printLog bool = true
	loggable := req.Context().Value(PrintLogKeyWord)
	if loggable != nil {
		flag := loggable.(string)
		if flag == "0" {
			printLog = false
		}
	}

	if req.Body != nil && printLog {

		m, err := ioutil.ReadAll(req.Body)
		if err == nil && len(m) > 0 {
			fmt.Printf("请求数据 %s \n", string(m))

			rc := ioutil.NopCloser(bytes.NewReader(m))
			req.Body = rc
		}
	}

	result, err1 := t.Transport.RoundTrip(req)

	if result != nil && result.Body != nil && printLog {

		m, err := ioutil.ReadAll(result.Body)
		if err == nil && len(m) > 0 {
			fmt.Printf("返回数据 %d %s", result.StatusCode,
				strings.Replace(string(m), "\n", "", -1))
			rc := ioutil.NopCloser(bytes.NewReader(m))
			result.Body = rc
		}

	}

	timeoutError := req.Context().Value(KlookMockTimeout)
	if timeoutError != nil {
		key := timeoutError.(string)
		// 连续3次timeout 第四次成功 测试需求
		var timeout bool = false

		if keyVal, ok := t.timeoutKey[key]; ok {
			if keyVal >= 2 {
				timeout = false
				t.timeoutKey[key] = 1
			} else {
				timeout = true
				t.timeoutKey[key] = keyVal + 1
			}
		} else {
			timeout = true
			t.timeoutKey[key] = 1
		}

		if timeout {
			// 制作超时情况
			result.Body = nil
			result.StatusCode = 504
			err1 = platformTimeoutError{}
		}
	}

	return result, err1
}

var HttpUtil httpUtil = httpUtil{}
