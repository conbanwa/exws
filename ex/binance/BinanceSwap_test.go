package binance

import (
	"net/http"
	"net/url"
	"qa3/wstrader"
	"qa3/wstrader/config"
	"qa3/wstrader/cons"
	"testing"
	"time"
)

var bs = NewBinanceSwap(&wstrader.APIConfig{
	Endpoint: "https://testnet.binancefuture.com",
	HttpClient: &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse("socks5://" + config.Proxy)
			},
		},
		Timeout: 10 * time.Second,
	},
	ApiKey:       "",
	ApiSecretKey: "",
})

func TestBinanceSwap_Ping(t *testing.T) {
	bs.Ping()
}
func TestBinanceSwap_GetFutureDepth(t *testing.T) {
	t.Log(bs.GetFutureDepth(cons.BTC_USDT, "", 1))
}
func TestBinanceSwap_GetFutureIndex(t *testing.T) {
	t.Log(bs.GetFutureIndex(cons.BTC_USDT))
}
func TestBinanceSwap_GetKlineRecords(t *testing.T) {
	kline, err := bs.GetKlineRecords("", cons.BTC_USDT, cons.KLINE_PERIOD_4H, 1, wstrader.OptionalParameter{"test": 0})
	t.Log(err, kline[0].Kline)
}
func TestBinanceSwap_GetTrades(t *testing.T) {
	t.Log(bs.GetTrades("", cons.BTC_USDT, 0))
}
func TestBinanceSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(bs.GetFutureUserinfo())
}
func TestBinanceSwap_PlaceFutureOrder(t *testing.T) {
	t.Log(bs.PlaceFutureOrder(cons.BTC_USDT, "", "8322", "0.01", cons.OPEN_BUY, 0, 0))
}
func TestBinanceSwap_PlaceFutureOrder2(t *testing.T) {
	t.Log(bs.PlaceFutureOrder(cons.BTC_USDT, "", "8322", "0.01", cons.OPEN_BUY, 1, 0))
}
func TestBinanceSwap_GetFutureOrder(t *testing.T) {
	t.Log(bs.GetFutureOrder("1431689723", cons.BTC_USDT, ""))
}
func TestBinanceSwap_FutureCancelOrder(t *testing.T) {
	t.Log(bs.FutureCancelOrder(cons.BTC_USDT, "", "1431554165"))
}
func TestBinanceSwap_GetFuturePosition(t *testing.T) {
	t.Log(bs.GetFuturePosition(cons.BTC_USDT, ""))
}
