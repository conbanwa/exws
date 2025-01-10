package binance

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"net/http"
	"testing"
)

var baDapi = NewBinanceFutures(&wstrader.APIConfig{
	HttpClient:   http.DefaultClient,
	ApiKey:       apiKey,
	ApiSecretKey: apiSecretkey,
})

func TestBinanceFutures_GetFutureDepth(t *testing.T) {
	skipKey(t)
	t.Log(baDapi.GetFutureDepth(cons.ETH_USD, cons.QUARTER_CONTRACT, 10))
}
func TestBinanceSwap_GetFutureTicker(t *testing.T) {
	skipKey(t)
	ticker, err := baDapi.GetFutureTicker(cons.LTC_USD, cons.SWAP_CONTRACT)
	t.Log(err)
	t.Logf("%+v", ticker)
}
func TestBinance_GetExchangeInfo(t *testing.T) {
	skipKey(t)
	baDapi.GetExchangeInfo()
}
func TestBinanceFutures_GetFutureUserinfo(t *testing.T) {
	skipKey(t)
	t.Log(baDapi.GetFutureUserinfo())
}
func TestBinanceFutures_PlaceFutureOrder(t *testing.T) {
	skipKey(t)
	//1044675677
	t.Log(baDapi.PlaceFutureOrder(cons.BTC_USD, cons.QUARTER_CONTRACT, "19990", "2", cons.OPEN_SELL, 0, 10))
}
func TestBinanceFutures_LimitFuturesOrder(t *testing.T) {
	skipKey(t)
	t.Log(baDapi.LimitFuturesOrder(cons.BTC_USD, cons.QUARTER_CONTRACT, "20001", "2", cons.OPEN_SELL))
}
func TestBinanceFutures_MarketFuturesOrder(t *testing.T) {
	skipKey(t)
	t.Log(baDapi.MarketFuturesOrder(cons.BTC_USD, cons.QUARTER_CONTRACT, "2", cons.OPEN_SELL))
}
func TestBinanceFutures_GetFutureOrder(t *testing.T) {
	skipKey(t)
	t.Log(baDapi.GetFutureOrder("1045208666", cons.BTC_USD, cons.QUARTER_CONTRACT))
}
func TestBinanceFutures_FutureCancelOrder(t *testing.T) {
	skipKey(t)
	t.Log(baDapi.FutureCancelOrder(cons.BTC_USD, cons.QUARTER_CONTRACT, "1045328328"))
}
func TestBinanceFutures_GetFuturePosition(t *testing.T) {
	skipKey(t)
	t.Log(baDapi.GetFuturePosition(cons.BTC_USD, cons.QUARTER_CONTRACT))
}
func TestBinanceFutures_GetUnfinishFutureOrders(t *testing.T) {
	skipKey(t)
	t.Log(baDapi.GetUnfinishFutureOrders(cons.BTC_USD, cons.QUARTER_CONTRACT))
}
