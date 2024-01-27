package config

import (
	"net/url"
	"os"
	"time"
)

var (
	Proxy                        = "127.0.0.1:7890"
	UseProxy                bool = false // = isWindows()
	DefaultHttpClientConfig      = &HttpClientConfig{
		Proxy:        nil,
		HttpTimeout:  2 * time.Second,
		MaxIdleConns: 10}
)

func SetProxy() {
	if !UseProxy {
		return
	}
	os.Setenv("HTTPS_PROXY", "socks5://"+Proxy)
	DefaultHttpClientConfig.Proxy = &url.URL{
		Scheme: "socks5",
		Host:   Proxy,
	}
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
