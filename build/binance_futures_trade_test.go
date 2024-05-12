package build

import (
	"github.com/conbanwa/logs"
	"qa3/wstrader"
	"qa3/wstrader/cons"
	"qa3/wstrader/ex/binance"
	"qa3/wstrader/q"
	"testing"
)

const (
	BinanceTestnetApiKey       = "YOUR_KEY"
	BinanceTestnetApiKeySecret = "YOUR_KEY_SECRET"
)

func TestFetchFutureDepthAndIndex(t *testing.T) {
	binanceApi := DefaultAPIBuilder.APIKey(BinanceTestnetApiKey).APISecretkey(BinanceTestnetApiKeySecret).Endpoint(binance.TestnetSpotWsBaseUrl).BuildFuture(cons.BINANCE_SWAP)
	depth, err := binanceApi.GetFutureDepth(cons.BTC_USD, cons.SWAP_USDT_CONTRACT, 100)
	if err != nil {
		logs.F(err.Error())
	}
	askTotalAmount, bidTotalAmount := 0.0, 0.0
	askTotalVol, bidTotalVol := 0.0, 0.0
	for _, v := range depth.AskList {
		askTotalAmount += v.Amount
		askTotalVol += v.Price * v.Amount
	}
	for _, v := range depth.BidList {
		bidTotalAmount += v.Amount
		bidTotalVol += v.Price * v.Amount
	}
	markPrice, err := binanceApi.GetFutureIndex(cons.BTC_USD)
	if err != nil {
		logs.F(err.Error())
	}
	logs.Infof("CURRENT mark price: %f", markPrice)
	logs.Infof("ContractType: %s ContractId: %s Pair: %s UTime: %s AmountTickSize: %d\n", depth.ContractType, depth.ContractId, depth.Pair, depth.UTime.String(), depth.Pair.AmountTickSize)
	logs.Infof("askTotalAmount: %f, bidTotalAmount: %f, askTotalVol: %f, bidTotalVol: %f", askTotalAmount, bidTotalAmount, askTotalVol, bidTotalVol)
	logs.Infof("ask price averge: %f, bid price averge: %f,", askTotalVol/askTotalAmount, bidTotalVol/bidTotalAmount)
	logs.Infof("ask-bid spread: %f%%,", 100*(depth.AskList[0].Price-depth.BidList[0].Price)/markPrice)
}
func TestSubscribeSpotMarketData(t *testing.T) {
	binanceWs, err := DefaultAPIBuilder.APIKey(BinanceTestnetApiKey).APISecretkey(BinanceTestnetApiKeySecret).Endpoint(binance.TestnetFutureUsdBaseUrl).BuildSpotWs(cons.BINANCE)
	if err != nil {
		logs.F(err.Error())
	}
	binanceWs.TickerCallback(func(ticker *wstrader.Ticker) {
		logs.Infof("%+v\n", *ticker)
	})
	binanceWs.SubscribeTicker(cons.BTC_USDT)
	binanceWs.DepthCallback(func(depth *wstrader.Depth) {
		logs.Infof("%+v\n", *depth)
	})
	binanceWs.SubscribeDepth(cons.BTC_USDT)
	binanceWs.TradeCallback(func(trade *q.Trade) {
		logs.Infof("%+v\n", *trade)
	})
	binanceWs.SubscribeTrade(cons.BTC_USDT)
	select {}
}

func TestSubscribeFutureMarketData(t *testing.T) {
	binanceWs, err := DefaultAPIBuilder.APIKey(BinanceTestnetApiKey).APISecretkey(BinanceTestnetApiKeySecret).Endpoint(binance.TestnetFutureUsdWsBaseUrl).BuildFuturesWs(cons.BINANCE_FUTURES)
	if err != nil {
		logs.F(err.Error())
	}
	binanceWs.TickerCallback(func(ticker *wstrader.FutureTicker) {
		//logs.Infof("%+v\n", *ticker.Ticker)
	})
	binanceWs.SubscribeTicker(cons.BTC_USD, cons.SWAP_USDT_CONTRACT)
	binanceWs.DepthCallback(func(depth *wstrader.Depth) {
		logs.Infof("%+v\n", *depth)
	})
	binanceWs.SubscribeDepth(cons.BTC_USDT, cons.SWAP_USDT_CONTRACT)
	binanceWs.TradeCallback(func(trade *q.Trade, contractType string) {
		logs.Infof("%+v\n", *trade)
	})
	binanceWs.SubscribeTrade(cons.BTC_USDT, cons.SWAP_USDT_CONTRACT)
	select {}
}
