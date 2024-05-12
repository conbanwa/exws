package zb

import (
	"net/http"
	"qa3/wstrader/cons"
	"testing"
)

var (
	apiKey       = ""
	apiSecretkey = ""
	zb           = New(http.DefaultClient, apiKey, apiSecretkey)
)

func TestZb_GetAccount(t *testing.T) {
	acc, err := zb.GetAccount()
	t.Log(err)
	t.Log(acc.SubAccounts[cons.BTC])
}
func TestZb_GetTicker(t *testing.T) {
	ticker, _ := zb.GetTicker(cons.BCH_USD)
	t.Log(ticker)
}
func TestZb_GetDepth(t *testing.T) {
	dep, _ := zb.GetDepth(2, cons.BCH_USDT)
	t.Log(dep)
}
func TestZb_LimitSell(t *testing.T) {
	ord, err := zb.LimitSell("0.001", "75000", cons.NewCurrencyPair2("BTC_QC"))
	t.Log(err)
	t.Log(ord)
}
func TestZb_LimitBuy(t *testing.T) {
	ord, err := zb.LimitBuy("2", "4", cons.NewCurrencyPair2("1ST_QC"))
	t.Log(err)
	t.Log(ord)
}
func TestZb_CancelOrder(t *testing.T) {
	r, err := zb.CancelOrder("201802014255365", cons.NewCurrencyPair2("BTC_QC"))
	t.Log(err)
	t.Log(r)
}
func TestZb_GetUnfinishOrders(t *testing.T) {
	ords, err := zb.GetUnfinishedOrders(cons.NewCurrencyPair2("1ST_QC"))
	t.Log(err)
	t.Log(ords)
}
func TestZb_GetOneOrder(t *testing.T) {
	ord, err := zb.GetOneOrder("20180201341043", cons.NewCurrencyPair2("1ST_QC"))
	t.Log(err)
	t.Log(ord)
}
