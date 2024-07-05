package coinbase

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"net/http"
	"testing"

	"github.com/conbanwa/logs"
)

var api = New(http.DefaultClient, "", "")

func TestCoinbase_GetTicker(t *testing.T) {
	ticker, err := api.GetTicker(cons.BTC_USD)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
func TestCoinbase_Get24HStats(t *testing.T) {
	stats, err := api.Get24HStats(cons.BTC_USD)
	t.Log("err=>", err)
	t.Log("stats=>", stats)
}
func TestCoinbase_GetDepth(t *testing.T) {
	dep, err := api.GetDepth(2, cons.BTC_USD)
	t.Log("err=>", err)
	t.Log("bids=>", dep.BidList)
	t.Log("asks=>", dep.AskList)
}
func TestCoinbase_GetKlineRecords(t *testing.T) {
	logs.Log.Level = logs.L_DEBUG
	t.Log(api.GetKlineRecords(cons.BTC_USD, cons.KLINE_PERIOD_1DAY, 0, wstrader.OptionalParameter{"test": 0}))
}
