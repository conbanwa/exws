package okx

import (
	"fmt"
	"net/url"
	"qa3/wstrader"
	"qa3/wstrader/cons"
	"testing"
)

func newOKExV5Client() *OKX {
	return NewOKExV5(&wstrader.APIConfig{
		//HttpClient: &http.Client{
		//	Transport: &http.Transport{
		//		Proxy: func(req *http.Request) (*url.URL, error) {
		//			return &url.URL{
		//				Scheme: "socks5",
		//				Host:   conf.PROXY}, nil
		//		},
		//	},
		//},
		Endpoint:      "https://www.okx.com",
		ApiKey:        "",
		ApiSecretKey:  "",
		ApiPassphrase: "",
	})
}
func TestOKExV5_GetTicker(t *testing.T) {
	o := newOKExV5Client()
	fmt.Println(o.GetTickerV5("BTC-USD-SWAP"))
}
func TestOKExV5_GetDepth(t *testing.T) {
	o := newOKExV5Client()
	fmt.Println(o.GetDepthV5("BTC-USD-SWAP", 0))
}
func TestOKExV5_GetKlineRecordsV5(t *testing.T) {
	o := newOKExV5Client()
	fmt.Println(o.GetKlineRecordsV5("BTC-USD-SWAP", cons.KLINE_PERIOD_1H, &url.Values{}))
}
