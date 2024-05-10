package okex

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/config"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"testing"
	"time"

	"github.com/conbanwa/logs"
)

func init() {
	logs.Log.Level = logs.L_DEBUG
}
func TestNewOKExSpotV3Ws(t *testing.T) {
	config.SetProxy()
	okexSpotV3Ws := okex.OKExV3SpotWs
	okexSpotV3Ws.TickerCallback(func(ticker *wstrader.Ticker) {
		t.Log(ticker)
	})
	okexSpotV3Ws.DepthCallback(func(depth *wstrader.Depth) {
		t.Log(depth)
	})
	okexSpotV3Ws.TradeCallback(func(trade *q.Trade) {
		t.Log(trade)
	})
	okexSpotV3Ws.KLineCallback(func(kline *wstrader.Kline, period cons.KlinePeriod) {
		t.Log(period, kline)
	})
	//okexSpotV3Ws.SubscribeDepth(module.EOS_USDT, 5)
	//okexSpotV3Ws.SubscribeTrade(module.EOS_USDT)
	//okexSpotV3Ws.SubscribeTicker(module.EOS_USDT)
	okexSpotV3Ws.SubscribeKline(cons.EOS_USDT, cons.KLINE_PERIOD_1H)
	time.Sleep(time.Minute)
}
