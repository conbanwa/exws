package gateio

import (
	"github.com/conbanwa/exws/cons"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var gate = New(http.DefaultClient, "", "")

func TestGate_AllTicker(t *testing.T) {
	ticker, err := gate.AllTicker(nil)
	assert.Nil(t, err)
	t.Logf("ticker=>%+v", ticker)
}
func TestGate_GetTicker(t *testing.T) {
	ticker, err := gate.GetTicker(cons.BTC_USDT)
	assert.Nil(t, err)
	t.Logf("ticker=>%+v", ticker)
}
func TestGate_GetDepth(t *testing.T) {
	dep, err := gate.GetDepth(1, cons.BTC_USDT)
	assert.Nil(t, err)
	t.Log("asks=>", dep.AskList[0], dep.AskList[1], dep.AskList[2])
	t.Log("bids=>", dep.BidList[0], dep.BidList[1], dep.BidList[2])
}
