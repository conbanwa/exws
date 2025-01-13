package web

import (
	"fmt"
	"github.com/conbanwa/exws/config"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	fastHttpClient = &fasthttp.Client{
		Name:                "goex-http-utils",
		MaxConnsPerHost:     16,
		MaxIdleConnDuration: 20 * time.Second,
		ReadTimeout:         10 * time.Second,
		WriteTimeout:        10 * time.Second,
	}
	socksDialer fasthttp.DialFunc
)

func init() {
	setProxy()
}

func setProxy() {
	url := config.GetProxy(true)
	if url == "" {
		return
	}
	socksDialer = fasthttpproxy.FasthttpSocksDialer(url)
	fastHttpClient.Dial = socksDialer
}

func FasthttpRequest(client *http.Client, reqMethod, reqUrl, postData string, headers map[string]string) (body []byte, err error) {
	if transport := client.Transport; transport != nil && config.UseProxy {
		if proxy, err := transport.(*http.Transport).Proxy(nil); err == nil && proxy != nil {
			if proxyUrl := proxy.String(); proxy.Scheme != "socks5" {
				panic("fasthttp only support the socks5 proxy " + proxy.Scheme + proxyUrl)
			} else if socksDialer == nil {
				setProxy()
			}
		}
	}
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
	}()
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	req.Header.SetMethod(reqMethod)
	req.SetRequestURI(reqUrl)
	req.SetBodyString(postData)
	err = fastHttpClient.Do(req, resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() == http.StatusTeapot || resp.StatusCode() == http.StatusTooManyRequests {
		panic("FATAL breaking a request rate limit")
	}
	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("HttpStatusCode:%d", resp.StatusCode())
	}
	body = resp.Body()
	return
}
func NewRequest(client *http.Client, reqType, reqUrl, postData string, requestHeaders map[string]string) (body []byte, err error) {
	// fasthttp := os.Getenv("HTTP_LIB")
	// if fasthttp == "fasthttp" {
	// 	logs.E(fasthttp)
	return FasthttpRequest(client, reqType, reqUrl, postData, requestHeaders)
	// }
}

func HttpRequest(client *http.Client, reqType, reqUrl, postData string, requestHeaders map[string]string) (body []byte, err error) {
	req, err := http.NewRequest(reqType, reqUrl, strings.NewReader(postData))
	if err != nil {
		return nil, err
	}
	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36")
	}
	for k, v := range requestHeaders {
		req.Header.Add(k, v)
	}
	log.Info().Any("Header", req.Header).Send()
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Any("req", *req).Any("response data", resp).Send()
		return
	}
	defer func(body io.ReadCloser) {
		err = body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Any("response data", resp.Body).Send()
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("HttpStatusCode:%d", resp.StatusCode)
	}
	return
}

func Get(url string) ([]byte, error) {
	return NewRequest(new(http.Client), http.MethodGet, url, "", nil)
}
