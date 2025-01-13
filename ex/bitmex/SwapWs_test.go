package bitmex

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/config"
	"github.com/conbanwa/exws/cons"
	"testing"
	"time"
)

func TestNewSwapWs(t *testing.T) {
	config.SetProxy()
	ws := NewSwapWs()
	ws.DepthCallback(func(depth *exws.Depth) {
		t.Log(depth)
	})
	ws.TickerCallback(func(ticker *exws.FutureTicker) {
		t.Logf("%s %v", ticker.ContractType, ticker.Ticker)
	})
	//ws.SubscribeDepth(cons.NewCurrencyPair2("LTC_USD"), cons.SWAP_CONTRACT)
	ws.SubscribeTicker(cons.LTC_USDT, cons.SWAP_CONTRACT)
	time.Sleep(5 * time.Second)
}
