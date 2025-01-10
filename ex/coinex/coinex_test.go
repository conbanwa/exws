package coinex

import (
	"fmt"
	"github.com/conbanwa/wstrader/cons"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const (
	TestKey    = "YOUR_KEY"
	TestSecret = "YOUR_KEY_SECRET"
)

func skipKey(t *testing.T) {
	if TestKey == "YOUR_KEY" {
		t.Skip("Skipping testing without TestKey")
	}
}

var coinex = New(http.DefaultClient, TestKey, TestSecret)

func TestCoinEx_GetTicker(t *testing.T) {
	ticker, err := coinex.GetTicker(cons.LTC_BTC)
	assert.Nil(t, err)
	t.Log(ticker)
}
func TestCoinEx_GetDepth(t *testing.T) {
	dep, err := coinex.GetDepth(5, cons.LTC_BTC)
	assert.Nil(t, err)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}
func TestCoinEx_GetAccount(t *testing.T) {
	skipKey(t)
	acc, err := coinex.GetAccount()
	assert.Nil(t, err)
	t.Log(acc)
}
func TestCoinEx_LimitBuy(t *testing.T) {
}
func TestCoinEx_LimitSell(t *testing.T) {
	skipKey(t)
	ord, err := coinex.LimitSell("100", "0.0000601", cons.NewCurrencyPair2("CET_BCH"))
	assert.Nil(t, err)
	t.Log(ord)
}
func TestCoinEx_GetUnfinishOrders(t *testing.T) {
	skipKey(t)
	ords, err := coinex.GetUnfinishedOrders(cons.NewCurrencyPair2("CET_BCH"))
	assert.Nil(t, err)
	if len(ords) > 0 {
		t.Log(fmt.Sprint(ords[0].OrderID))
	}
}
func TestCoinEx_CancelOrder(t *testing.T) {
	skipKey(t)
	r, err := coinex.CancelOrder("37504128", cons.NewCurrencyPair2("CET_BCH"))
	assert.Nil(t, err)
	t.Log(r)
}
func TestCoinEx_GetOneOrder(t *testing.T) {
	skipKey(t)
	ord, err := coinex.GetOneOrder("37504128", cons.NewCurrencyPair2("CET_BCH"))
	assert.Nil(t, err)
	t.Log(ord)
}
