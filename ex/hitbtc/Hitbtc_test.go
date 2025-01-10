package hitbtc

import (
	"github.com/conbanwa/wstrader/cons"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
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

var htb *Hitbtc

func init() {
	htb = New(http.DefaultClient, TestKey, TestSecret)
}
func TestHitbtc_GetSymbols(t *testing.T) {
	// panic: (map[string]interface {}) 0xc0000a1350 [recovered]
	// t.Log(htb.GetSymbols())
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
	skipKey(t)
	res, err := htb.GetAccount()
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestDepth(t *testing.T) {
	res, err := htb.GetDepth(10, cons.BTC_USD)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestKline(t *testing.T) {
	res, err := htb.GetKline(cons.BTC_USD, "1M", 10, 0)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestTrades(t *testing.T) {
	res, err := htb.GetTrades(cons.BTC_USD, 1519862400)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestPlaceOrder(t *testing.T) {
	skipKey(t)
	res, err := htb.LimitBuy("15", "0.000008", cons.BTC_USD)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestCancelOrder(t *testing.T) {
	skipKey(t)
	res, err := htb.CancelOrder("a605f2abbcc750da9138687bb27a2835", cons.BTC_USD)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestGetOneOrder(t *testing.T) {
	skipKey(t)
	res, err := htb.GetOneOrder("177836e71c8d57a14648d465e893efce", cons.BTC_USD)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestGetOrders(t *testing.T) {
	skipKey(t)
	res, err := htb.GetOrderHistorys(cons.BTC_USD)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
func TestGetUnfinishOrders(t *testing.T) {
	skipKey(t)
	res, err := htb.GetUnfinishedOrders(cons.BTC_USD)
	requires := require.New(t)
	requires.Nil(err)
	t.Log(res)
}
