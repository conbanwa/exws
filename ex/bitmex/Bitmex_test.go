package bitmex

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/config"
	"github.com/conbanwa/wstrader/cons"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var httpProxyClient = &http.Client{
	Transport: &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return &url.URL{
				Scheme: "socks5",
				Host:   config.Proxy}, nil
		},
	},
	Timeout: 10 * time.Second,
}

func init() {
	mex = New(&wstrader.APIConfig{
		Endpoint:   "https://testnet.bitmex.com/",
		HttpClient: httpProxyClient,
	})
}

var mex *Bitmex

func TestBitmex_GetFutureDepth(t *testing.T) {
	dep, err := mex.GetFutureDepth(cons.ETH_USDT, cons.SWAP_CONTRACT, 5)
	assert.Nil(t, err)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}
func TestBitmex_GetFutureTicker(t *testing.T) {
	ticker, err := mex.GetFutureTicker(cons.BTC_USD, "")
	if assert.Nil(t, err) {
		t.Logf("buy:%.8f ,sell: %.8f ,Last:%.8f , vol:%.8f", ticker.Buy, ticker.Sell, ticker.Last, ticker.Vol)
	}
}
func TestBitmex_GetIndicativeFundingRate(t *testing.T) {
	rate, time, err := mex.GetIndicativeFundingRate("XBTUSD")
	if assert.Nil(t, err) {
		t.Log(rate)
		t.Log(time.Local())
	}
}
func TestBitmex_GetFutureUserinfo(t *testing.T) {
	// userinfo, err := mex.GetFutureUserinfo()
	// if assert.Nil(t, err) {
	// 	t.Logf("%.8f", userinfo.FutureSubAccounts[cons.BTC].AccountRights)
	// 	t.Logf("%.8f", userinfo.FutureSubAccounts[cons.BTC].KeepDeposit)
	// 	t.Logf("%.8f", userinfo.FutureSubAccounts[cons.BTC].ProfitReal)
	// 	t.Logf("%.8f", userinfo.FutureSubAccounts[cons.BTC].ProfitUnreal)
	// }
}
func TestBitmex_GetFuturePosition(t *testing.T) {
	t.Log(mex.GetFuturePosition(cons.BTC_USD, ""))
}
func TestBitmex_PlaceFutureOrder(t *testing.T) {
	//{"orderID":"ae0436f4-9229-0be1-e9ea-45073a2a404a","clOrdID":"goexba0c770d9cea445eafb12b95fe220a0f"
	t.Log(mex.PlaceFutureOrder(cons.BTC_USD, cons.SWAP_CONTRACT, "9999", "2", cons.CLOSE_SELL, 0, 10))
}
func TestBitmex_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(mex.GetUnfinishFutureOrders(cons.BTC_USD, cons.SWAP_CONTRACT))
}
func TestBitmex_GetFutureOrder(t *testing.T) {
	t.Log(mex.GetFutureOrder("ae0436f4-9229-0be1-e9ea-45073a2a404a", cons.BTC_USD, cons.SWAP_CONTRACT))
}
func TestBitmex_FutureCancelOrder(t *testing.T) {
	t.Log(mex.FutureCancelOrder(cons.BTC_USD, cons.SWAP_CONTRACT, "goexfd6fd7694877448e8ae81a9cd7ecd89a"))
}
