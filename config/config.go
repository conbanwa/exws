package config

import (
	"net/url"
	"os"
	"time"
)

var (
	UseProxy                = false
	Scheme                  = "socks5"
	Proxy                   = "127.0.0.1:7890"
	DefaultHttpClientConfig = &HttpClientConfig{
		Proxy:        nil,
		HttpTimeout:  2 * time.Second,
		MaxIdleConns: 20,
	}
)

func GetProxy(withScheme ...bool) string {
	if UseProxy {
		if len(withScheme) > 0 {
			return Scheme + "://" + Proxy
		}
		return Proxy
	}
	return ""
}

func SetProxy() {
	os.Setenv("HTTPS_PROXY", GetProxy(true))
	if !UseProxy {
		DefaultHttpClientConfig.Proxy = nil
		return
	}
	DefaultHttpClientConfig.Proxy = &url.URL{
		Scheme: Scheme,
		Host:   GetProxy(),
	}
}

func IsWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
