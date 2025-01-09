package huobi

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"testing"
	"time"
)

func TestNewHbdmWs(t *testing.T) {
	ws := NewHbdmWs()
	ws.SetCallbacks(func(ticker *wstrader.FutureTicker) {
		t.Log(ticker.Ticker)
	}, func(depth *wstrader.Depth) {
		t.Log(">>>>>>>>>>>>>>>")
		t.Log(depth.ContractType, depth.Pair)
		t.Log(depth.BidList)
		t.Log(depth.AskList)
	}, func(trade *q.Trade, s string) {
		t.Log(s, trade)
	})
	t.Log(ws.SubscribeTicker(cons.BTC_USD, cons.QUARTER_CONTRACT))
	t.Log(ws.SubscribeDepth(cons.BTC_USD, cons.NEXT_WEEK_CONTRACT))
	t.Log(ws.SubscribeTrade(cons.LTC_USD, cons.THIS_WEEK_CONTRACT))
	time.Sleep(time.Second)
}
