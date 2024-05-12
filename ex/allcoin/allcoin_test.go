package allcoin

import (
	"net/http"
	"qa3/wstrader/cons"
	"testing"
)

var ac = New(http.DefaultClient, "", "")

func TestAllcoin_GetAccount(t *testing.T) {
	return
	t.Log(ac.GetAccount())
}
func TestAllcoin_GetUnfinishedOrders(t *testing.T) {
	return
	t.Log(ac.GetUnfinishedOrders(cons.ETH_BTC))
}
func TestAllcoin_GetTicker(t *testing.T) {
	return
	t.Log(ac.GetTicker(cons.ETH_BTC))
}
func TestAllcoin_GetDepth(t *testing.T) {
	return
	dep, _ := ac.GetDepth(1, cons.ETH_BTC)
	t.Log(dep)
}
func TestAllcoin_LimitBuy(t *testing.T) {
	t.Log(ac.LimitBuy("1", "0.07", cons.ETH_BTC))
}
