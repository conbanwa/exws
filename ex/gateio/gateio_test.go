package gateio

import (
	"net/http"
	"github.com/conbanwa/wstrader/cons"
	"testing"

	"github.com/conbanwa/logs"
)

var gate = New(http.DefaultClient, "", "")

func TestGate_AllTicker(t *testing.T) {
	ticker, err := gate.AllTicker(nil)
	logs.D(ticker)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
func TestGate_GetTicker(t *testing.T) {
	ticker, err := gate.GetTicker(cons.BTC_USDT)
	t.Log("err=>", err)
	t.Log("ticker=>", ticker)
}
func TestGate_GetDepth(t *testing.T) {
	dep, err := gate.GetDepth(1, cons.BTC_USDT)
	t.Log("err=>", err)
	t.Log("asks=>", dep.AskList)
	t.Log("bids=>", dep.BidList)
}
