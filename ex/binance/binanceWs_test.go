package binance

import (
	"qa3/wstrader"
	"qa3/wstrader/config"
	"qa3/wstrader/cons"
	"testing"
	"time"
)

var spotWs *SpotWs

func createSpotWs() {
	config.SetProxy()
	spotWs = NewSpotWs()
	spotWs.DepthCallback(func(depth *wstrader.Depth) {
		log.Println(depth)
	})
	spotWs.TickerCallback(func(ticker *wstrader.Ticker) {
		log.Println(ticker)
	})
}
func TestSpotWs_DepthCallback(t *testing.T) {
	createSpotWs()
	spotWs.SubscribeDepth(cons.BTC_USDT)
	spotWs.SubscribeTicker(cons.LTC_USDT)
	time.Sleep(11 * time.Minute)
}
func TestSpotWs_SubscribeTicker(t *testing.T) {
	createSpotWs()
	spotWs.SubscribeTicker(cons.LTC_USDT)
	time.Sleep(30 * time.Minute)
}
