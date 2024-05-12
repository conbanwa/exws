package bittrex

import (
	"net/http"
	"qa3/wstrader/cons"
	"testing"
)

var b = New(http.DefaultClient, "", "")

func TestBittrex_GetTicker(t *testing.T) {
	ticker, err := b.GetTicker(cons.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
func TestBittrex_GetDepth(t *testing.T) {
	dep, err := b.GetDepth(1, cons.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ask=>", dep.AskList)
	t.Log("bid=>", dep.BidList)
}
