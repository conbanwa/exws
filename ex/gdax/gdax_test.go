package gdax

import (
	"net/http"
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"testing"

	"github.com/conbanwa/logs"
)

var gdax = New(http.DefaultClient, "", "")

func TestGdax_GetTicker(t *testing.T) {
	ticker, err := gdax.GetTicker(cons.BTC_USD)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
func TestGdax_Get24HStats(t *testing.T) {
	stats, err := gdax.Get24HStats(cons.BTC_USD)
	t.Log("err=>", err)
	t.Log("stats=>", stats)
}
func TestGdax_GetDepth(t *testing.T) {
	dep, err := gdax.GetDepth(2, cons.BTC_USD)
	t.Log("err=>", err)
	t.Log("bids=>", dep.BidList)
	t.Log("asks=>", dep.AskList)
}
func TestGdax_GetKlineRecords(t *testing.T) {
	logs.Log.Level = logs.L_DEBUG
	t.Log(gdax.GetKlineRecords(cons.BTC_USD, cons.KLINE_PERIOD_1DAY, 0, wstrader.OptionalParameter{"test": 0}))
}
