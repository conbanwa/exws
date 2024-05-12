package okex

import (
	"net/http"
	"qa3/wstrader"
	"qa3/wstrader/config"
	"qa3/wstrader/cons"
	"qa3/wstrader/q"
	"testing"
	"time"

	"github.com/conbanwa/logs"
)

func init() {
	logs.Log.Level = logs.L_DEBUG
}
func TestNewOKExV3SwapWs(t *testing.T) {
	config.SetProxy()
	ok := NewOKEx(&wstrader.APIConfig{
		HttpClient: http.DefaultClient,
	})
	ok.OKExV3SwapWs.TickerCallback(func(ticker *wstrader.FutureTicker) {
		t.Log(ticker.Ticker, ticker.ContractType)
	})
	ok.OKExV3SwapWs.DepthCallback(func(depth *wstrader.Depth) {
		t.Log(depth)
	})
	ok.OKExV3SwapWs.TradeCallback(func(trade *q.Trade, s string) {
		t.Log(s, trade)
	})
	ok.OKExV3SwapWs.SubscribeTicker(cons.BTC_USDT, cons.SWAP_CONTRACT)
	time.Sleep(1 * time.Minute)
}
