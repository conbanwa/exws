package binance

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/config"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/stat/zelo"
	"github.com/stretchr/testify/assert"
	"testing"
)

var futuresWs *FuturesWs

func init() {
	log = zelo.Writer
	config.SetProxy()
	futuresWs = NewFuturesWs()
	futuresWs.DepthCallback(func(depth *wstrader.Depth) {
		//log.Debug().Any("depth", depth).Send()
	})
	futuresWs.TickerCallback(func(ticker *wstrader.FutureTicker) {
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
