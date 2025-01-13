package bitstamp

import (
	"github.com/conbanwa/exws/cons"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	apiKey       = ""
	apiSecretkey = "YOUR_KEY_SECRET"
)

var client = http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		log.Println("======")
		return nil
	},
}

var btmp = NewBitstamp(&client, apiKey, apiSecretkey, "")

func skipKey(t *testing.T) {
	if apiKey == "" {
		t.Skip("Skipping testing without apiKey")
	}
}

func TestBitstamp_GetAccount(t *testing.T) {
	skipKey(t)
	acc, err := btmp.GetAccount()
	assert.Nil(t, err)
	t.Log(acc)
}
func TestBitstamp_GetTicker(t *testing.T) {
	skipKey(t)
	ticker, err := btmp.GetTicker(cons.BTC_USD)
	assert.Nil(t, err)
	t.Log(ticker)
}
func TestBitstamp_GetDepth(t *testing.T) {
	skipKey(t)
	dep, err := btmp.GetDepth(5, cons.BTC_USD)
	assert.Nil(t, err)
	t.Log(dep.BidList)
	t.Log(dep.AskList)
}
func TestBitstamp_LimitBuy(t *testing.T) {
	skipKey(t)
	ord, err := btmp.LimitBuy("55", "0.12", cons.XRP_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
func TestBitstamp_LimitSell(t *testing.T) {
	skipKey(t)
	ord, err := btmp.LimitSell("40", "0.22", cons.XRP_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
func TestBitstamp_MarketBuy(t *testing.T) {
	skipKey(t)
	ord, err := btmp.MarketBuy("1", "", cons.XRP_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
func TestBitstamp_MarketSell(t *testing.T) {
	skipKey(t)
	ord, err := btmp.MarketSell("2", "", cons.XRP_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
func TestBitstamp_CancelOrder(t *testing.T) {
	skipKey(t)
	r, err := btmp.CancelOrder("311242779", cons.XRP_USD)
	assert.Nil(t, err)
	t.Log(r)
}
func TestBitstamp_GetUnfinishOrders(t *testing.T) {
	skipKey(t)
	ords, err := btmp.GetUnfinishedOrders(cons.XRP_USD)
	assert.Nil(t, err)
	t.Log(ords)
}
func TestBitstamp_GetOneOrder(t *testing.T) {
	skipKey(t)
	ord, err := btmp.GetOneOrder("311752078", cons.XRP_USD)
	assert.Nil(t, err)
	t.Log(ord)
}
