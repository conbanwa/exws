package hitbtc

import (
	"github.com/conbanwa/wstrader/cons"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	PubKey    = ""
	SecretKey = ""
)

var htb *Hitbtc

func init() {
	htb = New(http.DefaultClient, PubKey, SecretKey)
}
func TestHitbtc_GetSymbols(t *testing.T) {
	t.Log(htb.GetSymbols())
}
func TestHitbtc_adaptSymbolToCurrencyPair(t *testing.T) {
	t.Log(htb.adaptSymbolToCurrencyPair("DOGEBTC").String() == "DOGE_BTC")
	t.Log(htb.adaptSymbolToCurrencyPair("BTCGUSD").String() == "BTC_GUSD")
	t.Log(htb.adaptSymbolToCurrencyPair("btctusd").String() == "BTC_TUSD")
	t.Log(htb.adaptSymbolToCurrencyPair("BTCUSDC").String() == "BTC_USDC")
	t.Log(htb.adaptSymbolToCurrencyPair("ETHEOS").String() == "ETH_EOS")
}
func TestGetTicker(t *testing.T) {
	res, err := htb.GetTicker(cons.BCH_USD)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestGetAccount(t *testing.T) {
	res, err := htb.GetAccount()
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestDepth(t *testing.T) {
	res, err := htb.GetDepth(10, YCC_BTC)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestKline(t *testing.T) {
	res, err := htb.GetKline(YCC_BTC, "1M", 10, 0)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestTrades(t *testing.T) {
	res, err := htb.GetTrades(YCC_BTC, 1519862400)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestPlaceOrder(t *testing.T) {
	res, err := htb.LimitBuy("15", "0.000008", YCC_BTC)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestCancelOrder(t *testing.T) {
	res, err := htb.CancelOrder("a605f2abbcc750da9138687bb27a2835", YCC_BTC)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestGetOneOrder(t *testing.T) {
	res, err := htb.GetOneOrder("177836e71c8d57a14648d465e893efce", YCC_BTC)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestGetOrders(t *testing.T) {
	res, err := htb.GetOrderHistorys(YCC_BTC)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestGetUnfinishOrders(t *testing.T) {
	res, err := htb.GetUnfinishedOrders(YCC_BTC)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
