package binance

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/config"
	"github.com/conbanwa/exws/cons"
	"github.com/stretchr/testify/assert"
	"testing"
)

var spotWs *SpotWs

func createSpotWs() {
	config.SetProxy()
	spotWs = NewSpotWs()
	spotWs.DepthCallback(func(depth *exws.Depth) {
		log.Println(depth)
	})
	spotWs.TickerCallback(func(ticker *exws.Ticker) {
		log.Println(ticker)
	})
}
func TestSpotWs_DepthCallback(t *testing.T) {
	createSpotWs()
	assert.Nil(t, spotWs.SubscribeDepth(cons.BTC_USDT))
	assert.Nil(t, spotWs.SubscribeTicker(cons.LTC_USDT))
}
func TestSpotWs_SubscribeTicker(t *testing.T) {
	createSpotWs()
	assert.Nil(t, spotWs.SubscribeTicker(cons.LTC_USDT))
}
