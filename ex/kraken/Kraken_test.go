package kraken

import (
	"net/http"
	"github.com/conbanwa/wstrader/cons"
	"testing"

	"github.com/stretchr/testify/assert"
)

var k = New(http.DefaultClient, "", "")
var BCH_XBT = cons.NewCurrencyPair(cons.BCH, cons.XBT)

func TestKraken_GetDepth(t *testing.T) {
	dep, err := k.GetDepth(2, cons.BTC_USD)
	assert.Nil(t, err)
	t.Log(dep)
}
func TestKraken_GetTicker(t *testing.T) {
	ticker, err := k.GetTicker(cons.ETC_BTC)
	assert.Nil(t, err)
	t.Log(ticker)
}
func TestKraken_GetAccount(t *testing.T) {
	acc, err := k.GetAccount()
	assert.Nil(t, err)
	t.Log(acc)
}
func TestKraken_LimitSell(t *testing.T) {
	ord, err := k.LimitSell("0.01", "6900", cons.BTC_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
func TestKraken_LimitBuy(t *testing.T) {
	ord, err := k.LimitBuy("0.01", "6100", cons.NewCurrencyPair(cons.XBT, cons.USD))
	assert.Nil(t, err)
	t.Log(ord)
}
func TestKraken_GetUnfinishOrders(t *testing.T) {
	ords, err := k.GetUnfinishedOrders(cons.NewCurrencyPair(cons.XBT, cons.USD))
	assert.Nil(t, err)
	t.Log(ords)
}
func TestKraken_CancelOrder(t *testing.T) {
	r, err := k.CancelOrder("O6EAJC-YAC3C-XDEEXQ", cons.NewCurrencyPair(cons.XBT, cons.USD))
	assert.Nil(t, err)
	t.Log(r)
}
func TestKraken_GetTradeBalance(t *testing.T) {
	//	k.GetTradeBalance()
}
func TestKraken_GetOneOrder(t *testing.T) {
	ord, err := k.GetOneOrder("ODCRMQ-RDEID-CY334C", cons.BTC_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
