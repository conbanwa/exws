package okex

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/config"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"net/http"
	"testing"
	"time"

	"github.com/conbanwa/logs"
)

var (
	client *http.Client
)

func init() {
	logs.Log.Level = logs.L_DEBUG
}
func TestNewOKExV3FuturesWs(t *testing.T) {
	config.SetProxy()
	ok := NewOKEx(&wstrader.APIConfig{
		HttpClient: http.DefaultClient,
	})
	ok.OKExV3FuturesWs.TickerCallback(func(ticker *wstrader.FutureTicker) {
		t.Log(ticker.Ticker, ticker.ContractType)
	})
	ok.OKExV3FuturesWs.DepthCallback(func(depth *wstrader.Depth) {
		t.Log(depth)
	})
	ok.OKExV3FuturesWs.TradeCallback(func(trade *q.Trade, s string) {
		t.Log(s, trade)
	})
	//ok.OKExV3FuturesWs.SubscribeTicker(module.EOS_USD, module.QUARTER_CONTRACT)
	ok.OKExV3FuturesWs.SubscribeDepth(cons.EOS_USD, cons.QUARTER_CONTRACT)
	//ok.OKExV3FuturesWs.SubscribeTrade(module.EOS_USD, module.QUARTER_CONTRACT)
	time.Sleep(time.Second * 20)
}
