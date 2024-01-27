package coinbene

import (
	"net/http"
	"net/url"
	"qa3/wstrader"
	"qa3/wstrader/config"
	"qa3/wstrader/cons"
	"testing"
	"time"
)

var (
	httpProxyClient = &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "socks5",
					Host:   config.Proxy}, nil
			},
		},
		Timeout: 10 * time.Second,
	}
	coinbeneSwap = NewCoinbeneSwap(wstrader.APIConfig{
		HttpClient:   httpProxyClient,
		Endpoint:     "",
		ApiKey:       "",
		ApiSecretKey: "",
	})
)

func TestCoinbeneSwap_GetFutureTicker(t *testing.T) {
	t.Log(coinbeneSwap.GetFutureTicker(cons.BTC_USD, cons.SWAP_CONTRACT))
}
func TestCoinbeneSwap_GetFutureDepth(t *testing.T) {
	t.Log(coinbeneSwap.GetFutureDepth(cons.BTC_USDT, cons.SWAP_CONTRACT, 2))
}
func TestCoinbeneSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(coinbeneSwap.GetFutureUserinfo())
}
func TestCoinbeneSwap_GetFuturePosition(t *testing.T) {
	t.Log(coinbeneSwap.GetFuturePosition(cons.BTC_USDT, cons.SWAP_CONTRACT))
}
func TestCoinbeneSwap_PlaceFutureOrder(t *testing.T) {
	t.Log(coinbeneSwap.PlaceFutureOrder(cons.BTC_USDT, cons.SWAP_CONTRACT, "10000", "1", cons.OPEN_BUY, 0, 10))
}
func TestCoinbeneSwap_FutureCancelOrder(t *testing.T) {
	t.Log(coinbeneSwap.FutureCancelOrder(cons.BTC_USDT, cons.SWAP_CONTRACT, "580719990266232832"))
}
func TestCoinbeneSwap_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(coinbeneSwap.GetUnfinishFutureOrders(cons.BTC_USDT, cons.SWAP_CONTRACT))
}
func TestCoinbeneSwap_GetFutureOrder(t *testing.T) {
	t.Log(coinbeneSwap.GetFutureOrder("123", cons.BTC_USDT, cons.SWAP_CONTRACT))
}
