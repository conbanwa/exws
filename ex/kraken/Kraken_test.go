package kraken

import (
	"github.com/conbanwa/exws/cons"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	apiKey       = ""
	apiSecretkey = ""
)

var k = New(http.DefaultClient, apiKey, apiSecretkey)
var BCH_XBT = cons.NewCurrencyPair(cons.BCH, cons.XBT)

func skipKey(t *testing.T) {
	if apiKey == "" {
		t.Skip("Skipping testing without apiKey")
	}
}
func TestKraken_GetDepth(t *testing.T) {
	skipKey(t)
	dep, err := k.GetDepth(2, cons.BTC_USD)
	assert.Nil(t, err)
	t.Log(dep)
}
func TestKraken_GetTicker(t *testing.T) {
	skipKey(t)
	ticker, err := k.GetTicker(cons.ETC_BTC)
	assert.Nil(t, err)
	t.Log(ticker)
}
func TestKraken_GetAccount(t *testing.T) {
	skipKey(t)
	acc, err := k.GetAccount()
	assert.Nil(t, err)
	t.Log(acc)
}
func TestKraken_LimitSell(t *testing.T) {
	skipKey(t)
	ord, err := k.LimitSell("0.01", "6900", cons.BTC_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
func TestKraken_LimitBuy(t *testing.T) {
	skipKey(t)
	ord, err := k.LimitBuy("0.01", "6100", cons.NewCurrencyPair(cons.XBT, cons.USD))
	assert.Nil(t, err)
	t.Log(ord)
}
func TestKraken_GetUnfinishOrders(t *testing.T) {
	skipKey(t)
	ords, err := k.GetUnfinishedOrders(cons.NewCurrencyPair(cons.XBT, cons.USD))
	assert.Nil(t, err)
	t.Log(ords)
}
func TestKraken_CancelOrder(t *testing.T) {
	skipKey(t)
	r, err := k.CancelOrder("O6EAJC-YAC3C-XDEEXQ", cons.NewCurrencyPair(cons.XBT, cons.USD))
	assert.Nil(t, err)
	t.Log(r)
}
func TestKraken_GetTradeBalance(t *testing.T) {
	//	k.GetTradeBalance()
}
func TestKraken_GetOneOrder(t *testing.T) {
	skipKey(t)
	ord, err := k.GetOneOrder("ODCRMQ-RDEID-CY334C", cons.BTC_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
