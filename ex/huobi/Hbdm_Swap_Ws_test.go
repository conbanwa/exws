package huobi

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewHbdmSwapWs(t *testing.T) {
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
	assert.Nil(t, ws.SubscribeDepth(cons.BTC_USD, cons.SWAP_CONTRACT))
	assert.Nil(t, ws.SubscribeTicker(cons.BTC_USD, cons.SWAP_CONTRACT))
	assert.Nil(t, ws.SubscribeTrade(cons.BTC_USD, cons.SWAP_CONTRACT))
	time.Sleep(time.Second * 20)
}
