package binance

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/config"
	"github.com/conbanwa/exws/cons"
	"github.com/conbanwa/exws/stat/zelo"
	"github.com/stretchr/testify/assert"
	"testing"
)

var futuresWs *FuturesWs

func init() {
	log = zelo.Writer
	config.SetProxy()
	futuresWs = NewFuturesWs()
	futuresWs.DepthCallback(func(depth *exws.Depth) {
		//log.Debug().Any("depth", depth).Send()
	})
	futuresWs.TickerCallback(func(ticker *exws.FutureTicker) {
		//log.Println(ticker.Ticker, ticker.ContractType)
	})
}
func TestFuturesWs_DepthCallback(t *testing.T) {
	assert.Nil(t, futuresWs.SubscribeDepth(cons.LTC_USDT, cons.SWAP_USDT_CONTRACT))
	assert.Nil(t, futuresWs.SubscribeDepth(cons.LTC_USDT, cons.SWAP_CONTRACT))
	assert.Nil(t, futuresWs.SubscribeDepth(cons.LTC_USDT, cons.QUARTER_CONTRACT))
}
func TestFuturesWs_SubscribeTicker(t *testing.T) {
	assert.Nil(t, futuresWs.SubscribeTicker(cons.BTC_USDT, cons.SWAP_USDT_CONTRACT))
	assert.Nil(t, futuresWs.SubscribeTicker(cons.BTC_USDT, cons.SWAP_CONTRACT))
	assert.Nil(t, futuresWs.SubscribeTicker(cons.BTC_USDT, cons.QUARTER_CONTRACT))
}
