package huobi

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"testing"
	"time"

	"github.com/conbanwa/logs"
)

func TestNewHbdmSwapWs(t *testing.T) {
	logs.Log.Level = logs.L_DEBUG
	ws := NewHbdmSwapWs()
	ws.DepthCallback(func(depth *wstrader.Depth) {
		t.Log(depth)
	})
	ws.TickerCallback(func(ticker *wstrader.FutureTicker) {
		t.Log(ticker.Date, ticker.Last, ticker.Buy, ticker.Sell, ticker.High, ticker.Low, ticker.Vol)
	})
	ws.TradeCallback(func(trade *q.Trade, contract string) {
		t.Log(trade, contract)
	})
	//t.Log(ws.SubscribeDepth(module.BTC_USD, module.SWAP_CONTRACT))
	//t.Log(ws.SubscribeTicker(module.BTC_USD, module.SWAP_CONTRACT))
	t.Log(ws.SubscribeTrade(cons.BTC_USD, cons.SWAP_CONTRACT))
	time.Sleep(time.Minute)
}
