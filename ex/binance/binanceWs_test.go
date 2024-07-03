package binance

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/config"
	"github.com/conbanwa/wstrader/cons"
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
	time.Sleep(5 * time.Second)
}
func TestSpotWs_SubscribeTicker(t *testing.T) {
	createSpotWs()
	spotWs.SubscribeTicker(cons.LTC_USDT)
	time.Sleep(3 * time.Second)
}
