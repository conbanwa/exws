package bithumb

import (
	"github.com/conbanwa/wstrader/cons"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var bh = New(http.DefaultClient, "", "")

func TestBithumb_GetTicker(t *testing.T) {
	ticker, err := bh.GetTicker(cons.NewCurrencyPair2("ALL_KAW"))
	assert.Nil(t, err)
	t.Log("ticker=>", ticker)
}
func TestBithumb_GetDepth(t *testing.T) {
	dep, err := bh.GetDepth(1, cons.BTC_KRW)
	if assert.Nil(t, err) {
		t.Log("asks=>", dep.AskList)
		t.Log("bids=>", dep.BidList)
	}
}
