package bitfinex

import (
	"github.com/conbanwa/wstrader/cons"
	"net/http"
	"testing"
)

var bfx = New(http.DefaultClient, "", "")

func TestBitfinex_GetTicker(t *testing.T) {
	ticker, _ := bfx.GetTicker(cons.ETH_BTC)
	t.Log(ticker)
}
func TestBitfinex_GetDepth(t *testing.T) {
	dep, _ := bfx.GetDepth(2, cons.ETH_BTC)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}
func TestBitfinex_GetKline(t *testing.T) {
	kline, _ := bfx.GetKlineRecords(cons.BTC_USD, cons.KLINE_PERIOD_1MONTH, 10)
	for _, k := range kline {
		t.Log(k)
	}
}
