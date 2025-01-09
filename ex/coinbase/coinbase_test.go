package coinbase

import (
	"github.com/conbanwa/slice"
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"net/http"
	"testing"
)

var api = New(http.DefaultClient, "", "")

func TestCoinbase_GetTicker(t *testing.T) {
	ticker, err := api.GetTicker(cons.BTC_USD)
	if err != nil {
		t.Error(err)
	}
	t.Log("ticker=>", ticker)
}
func TestCoinbase_Get24HStats(t *testing.T) {
	stats, err := api.Get24HStats(cons.BTC_USD)
	if err != nil {
		t.Error(err)
	}
	t.Log("stats=>", stats)
}
func TestCoinbase_GetDepth(t *testing.T) {
	dep, err := api.GetDepth(2, cons.BTC_USD)
	if err != nil {
		t.Error(err)
	}
	t.Log("bids=>", slice.Slice(dep.BidList, 0, 4))
	t.Log("asks=>", slice.Slice(dep.AskList, -4))
}
func TestCoinbase_GetKlineRecords(t *testing.T) {
	kline, err := api.GetKlineRecords(cons.BTC_USD, cons.KLINE_PERIOD_1DAY, 0, wstrader.OptionalParameter{"test": 0})
	if err != nil {
		t.Error(err)
	}
	t.Log(slice.Slice(kline, 0, 4))
}
