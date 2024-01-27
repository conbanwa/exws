package bitget

import (
	"net/http"
	"net/url"
	"qa3/wstrader"
	"qa3/wstrader/config"
	"qa3/wstrader/cons"
	"testing"
)

var bg = NewSwap(&wstrader.APIConfig{
	HttpClient: &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "socks5",
					Host:   config.Proxy}, nil
			},
		},
	}, //需要代理的这样配置
	Endpoint: "https://capi.bitget.io",
	ClientId: "",
	Lever:    0,
})

func TestBitgetSwap_GetFutureTicker(t *testing.T) {
	t.Log(bg.GetFutureTicker(cons.ETH_USDT, ""))
}
func TestBitgetSwap_GetServerTime(t *testing.T) {
	t.Log(bg.GetServerTime())
}
func TestBitgetSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(bg.GetFutureUserinfo(cons.ETH_USDT))
}
func TestBitgetSwap_LimitFuturesOrder(t *testing.T) {
	t.Log(bg.LimitFuturesOrder(cons.ETH_USDT, "", "350", "1", cons.CLOSE_BUY))
}
func TestBitgetSwap_GetFuturePosition(t *testing.T) {
	t.Log(bg.GetFuturePosition(cons.ETH_USDT, ""))
}
func TestBitgetSwap_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(bg.GetUnfinishFutureOrders(cons.ETH_USDT, ""))
}
func TestBitgetSwap_SetMarginLevel(t *testing.T) {
	t.Log(bg.SetMarginLevel(cons.ETH_USDT, 10, 2))
}
func TestBitgetSwap_GetMarginLevel(t *testing.T) {
	t.Log(bg.GetMarginLevel(cons.ETH_USDT))
}
func TestBitgetSwap_GetContractInfo(t *testing.T) {
	t.Log(bg.GetContractInfo(cons.ETH_USDT))
}
func TestBitgetSwap_GetFutureOrder(t *testing.T) {
	t.Log(bg.GetFutureOrder("671529783552638913", cons.ETH_USDT, ""))
}
func TestBitgetSwap_FutureCancelOrder(t *testing.T) {
	t.Log(bg.FutureCancelOrder(cons.ETH_USDT, "", "671529783552638913"))
}
func TestBitgetSwap_ModifyAutoAppendMargin(t *testing.T) {
	t.Log(bg.ModifyAutoAppendMargin(cons.ETH_USDT, 1, 1))
}
