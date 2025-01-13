package coinbase

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/cons"
	"github.com/conbanwa/slice"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var api = New(http.DefaultClient, "", "")

func TestCoinbase_GetTicker(t *testing.T) {
	ticker, err := api.GetTicker(cons.BTC_USD)
	assert.Nil(t, err)
	t.Log("ticker=>", ticker)
}
func TestCoinbase_Get24HStats(t *testing.T) {
	stats, err := api.Get24HStats(cons.BTC_USD)
	assert.Nil(t, err)
	t.Log("stats=>", stats)
}
func TestCoinbase_GetDepth(t *testing.T) {
	dep, err := api.GetDepth(2, cons.BTC_USD)
	assert.Nil(t, err)
	t.Log("bids=>", slice.Slice(dep.BidList, 0, 4))
	t.Log("asks=>", slice.Slice(dep.AskList, -4))
}
func TestCoinbase_GetKlineRecords(t *testing.T) {
	kline, err := api.GetKlineRecords(cons.BTC_USD, cons.KLINE_PERIOD_1DAY, 0, exws.OptionalParameter{"test": 0})
	assert.Nil(t, err)
	t.Log(slice.Slice(kline, 0, 4))
}
