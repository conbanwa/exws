package huobi

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/config"
	"github.com/conbanwa/exws/cons"
	"testing"
	"time"
)

func TestNewSpotWs(t *testing.T) {
	config.SetProxy()
	spotWs := NewSpotWs()
	spotWs.DepthCallback(func(depth *exws.Depth) {
		t.Log("asks=", depth.AskList)
		t.Log("bids=", depth.BidList)
	})
	spotWs.TickerCallback(func(ticker *exws.Ticker) {
		t.Log(ticker)
	})
	spotWs.SubscribeTicker(cons.NewCurrencyPair2("BTC_USDT"))
	spotWs.SubscribeTicker(cons.NewCurrencyPair2("USDT_HUSD"))
	spotWs.SubscribeTicker(cons.NewCurrencyPair2("LTC_BTC"))
	spotWs.SubscribeTicker(cons.NewCurrencyPair2("EOS_ETH"))
	spotWs.SubscribeTicker(cons.NewCurrencyPair2("LTC_HT"))
	spotWs.SubscribeTicker(cons.NewCurrencyPair2("BTT_TRX"))
	spotWs.SubscribeDepth(cons.BTC_USDT)
	time.Sleep(time.Second * 3)
}
