package binance

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/config"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/stat/zelo"
	"os"
	"testing"
	"time"
)

var futuresWs *FuturesWs

func init() {
	log = zelo.Colored(os.Stderr)
	config.SetProxy()
	futuresWs = NewFuturesWs()
	futuresWs.DepthCallback(func(depth *wstrader.Depth) {
		log.Debug().Any("depth", depth).Send()
	})
	futuresWs.TickerCallback(func(ticker *wstrader.FutureTicker) {
		log.Println(ticker.Ticker, ticker.ContractType)
	})
}
func TestFuturesWs_DepthCallback(t *testing.T) {
	futuresWs.SubscribeDepth(cons.LTC_USDT, cons.SWAP_USDT_CONTRACT)
	futuresWs.SubscribeDepth(cons.LTC_USDT, cons.SWAP_CONTRACT)
	// futuresWs.SubscribeDepth(cons.LTC_USDT, cons.QUARTER_CONTRACT)
	time.Sleep(1 * time.Second)
}
func TestFuturesWs_SubscribeTicker(t *testing.T) {
	futuresWs.SubscribeTicker(cons.BTC_USDT, cons.SWAP_USDT_CONTRACT)
	futuresWs.SubscribeTicker(cons.BTC_USDT, cons.SWAP_CONTRACT)
	futuresWs.SubscribeTicker(cons.BTC_USDT, cons.QUARTER_CONTRACT)
	time.Sleep(3 * time.Second)
}
