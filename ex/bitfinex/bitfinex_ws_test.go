package bitfinex

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/cons"
	"github.com/conbanwa/exws/q"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestNewBitfinexWs(t *testing.T) {
	bitfinexWs := NewWs()
	handleBbo := func(ticker *q.Bbo) {
		log.Printf("Ticker: %+v: ", ticker)
	}
	handleTicker := func(ticker *exws.Ticker) {
		log.Printf("Ticker: %+v: ", ticker)
	}
	handleTrade := func(trade *q.Trade) {
		log.Printf("Trade: %+v: ", trade)
	}
	handleCandle := func(candle *exws.Kline) {
		log.Printf("Candle: %+v: ", candle)
	}
	bitfinexWs.SetCallbacks(handleBbo, handleTicker, handleTrade, handleCandle)
	//Ticker
	assert.Nil(t, bitfinexWs.SubscribeTicker(cons.BTC_USD))
	assert.Nil(t, bitfinexWs.SubscribeTicker(cons.LTC_USD))
	//Trades
	assert.Nil(t, bitfinexWs.SubscribeTrade(cons.BTC_USD))
	//Candles
	assert.Nil(t, bitfinexWs.SubscribeCandle(cons.BTC_USD, cons.KLINE_PERIOD_1MIN))
	time.Sleep(time.Second * 10)
}
