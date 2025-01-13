package okex

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/config"
	"github.com/conbanwa/exws/cons"
	"github.com/conbanwa/exws/q"
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
	ok := NewOKEx(&exws.APIConfig{
		HttpClient: http.DefaultClient,
	})
	ok.OKExV3FuturesWs.TickerCallback(func(ticker *exws.FutureTicker) {
		t.Log(ticker.Ticker, ticker.ContractType)
	})
	ok.OKExV3FuturesWs.DepthCallback(func(depth *exws.Depth) {
		t.Log(depth)
	})
	ok.OKExV3FuturesWs.TradeCallback(func(trade *q.Trade, s string) {
		t.Log(s, trade)
	})
	//ok.OKExV3FuturesWs.SubscribeTicker(cons.EOS_USD, cons.QUARTER_CONTRACT)
	ok.OKExV3FuturesWs.SubscribeDepth(cons.EOS_USD, cons.QUARTER_CONTRACT)
	//ok.OKExV3FuturesWs.SubscribeTrade(cons.EOS_USD, cons.QUARTER_CONTRACT)
	time.Sleep(time.Second * 10)
}
