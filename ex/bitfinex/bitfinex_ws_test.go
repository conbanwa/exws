package bitfinex

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"log"
	"testing"
	"time"
)

func TestNewBitfinexWs(t *testing.T) {
	bitfinexWs := NewWs()
	handleBbo := func(ticker *q.Bbo) {
		log.Printf("Ticker: %+v: ", ticker)
	}
	handleTicker := func(ticker *wstrader.Ticker) {
		log.Printf("Ticker: %+v: ", ticker)
	}
	handleTrade := func(trade *q.Trade) {
		log.Printf("Trade: %+v: ", trade)
	}
	handleCandle := func(candle *wstrader.Kline) {
		log.Printf("Candle: %+v: ", candle)
	}
	bitfinexWs.SetCallbacks(handleBbo, handleTicker, handleTrade, handleCandle)
	//Ticker
	t.Log(bitfinexWs.SubscribeTicker(cons.BTC_USD))
	t.Log(bitfinexWs.SubscribeTicker(cons.LTC_USD))
	//Trades
	t.Log(bitfinexWs.SubscribeTrade(cons.BTC_USD))
	//Candles
	t.Log(bitfinexWs.SubscribeCandle(cons.BTC_USD, cons.KLINE_PERIOD_1MIN))
	time.Sleep(time.Second * 10)
}
