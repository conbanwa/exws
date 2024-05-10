package binance

import (
	"qa3/wstrader"
	"qa3/wstrader/config"
	"qa3/wstrader/cons"
	"testing"
	"time"
)

var futuresWs *FuturesWs

func createFuturesWs() {
	config.SetProxy()
	futuresWs = NewFuturesWs()
	futuresWs.DepthCallback(func(depth *wstrader.Depth) {
		log.Println(depth)
	})
	futuresWs.TickerCallback(func(ticker *wstrader.FutureTicker) {
		log.Println(ticker.Ticker, ticker.ContractType)
	})
}
func TestFuturesWs_DepthCallback(t *testing.T) {
	createFuturesWs()
	futuresWs.SubscribeDepth(cons.LTC_USDT, cons.SWAP_USDT_CONTRACT)
	futuresWs.SubscribeDepth(cons.LTC_USDT, cons.SWAP_CONTRACT)
	futuresWs.SubscribeDepth(cons.LTC_USDT, cons.QUARTER_CONTRACT)
	time.Sleep(30 * time.Second)
}
func TestFuturesWs_SubscribeTicker(t *testing.T) {
	createFuturesWs()
	futuresWs.SubscribeTicker(cons.BTC_USDT, cons.SWAP_USDT_CONTRACT)
	futuresWs.SubscribeTicker(cons.BTC_USDT, cons.SWAP_CONTRACT)
	futuresWs.SubscribeTicker(cons.BTC_USDT, cons.QUARTER_CONTRACT)
	time.Sleep(30 * time.Second)
}
