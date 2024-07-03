package bitmex

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/config"
	"github.com/conbanwa/wstrader/cons"
	"testing"
	"time"
)

func TestNewSwapWs(t *testing.T) {
	config.SetProxy()
	ws := NewSwapWs()
	ws.DepthCallback(func(depth *wstrader.Depth) {
		t.Log(depth)
	})
	ws.TickerCallback(func(ticker *wstrader.FutureTicker) {
		t.Logf("%s %v", ticker.ContractType, ticker.Ticker)
	})
	//ws.SubscribeDepth(module.NewCurrencyPair2("LTC_USD"), module.SWAP_CONTRACT)
	ws.SubscribeTicker(cons.LTC_USDT, cons.SWAP_CONTRACT)
	time.Sleep(5 * time.Second)
}
