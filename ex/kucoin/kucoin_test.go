package kucoin

import (
	"github.com/conbanwa/wstrader/cons"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	apiKey    = ""
	apiSecretkey = "YOUR_KEY_SECRET"
)

func skipKey(t *testing.T) {
	if apiKey == "" {
		t.Skip("Skipping testing without apiKey")
	}
}

var kc = New(apiKey, apiSecretkey, "")

func TestKuCoinerrGetTicker(t *testing.T) {
	ticker, err := kc.GetTicker(cons.BTC_USDT)
	assert.Nil(t, err)
	t.Log(ticker)
}
func TestKuCoinerrGetDepth(t *testing.T) {
	depth, err := kc.GetDepth(10, cons.BTC_USDT)
	assert.Nil(t, err)
	t.Log(depth)
}
func TestKuCoinerrGetKlineRecords(t *testing.T) {
	kLines, err := kc.GetKlineRecords(cons.BTC_USDT, cons.KLINE_PERIOD_1MIN, 10)
	assert.Nil(t, err)
	t.Log(kLines)
}
func TestKuCoinerrGetTrades(t *testing.T) {
	trades, err := kc.GetTrades(cons.BTC_USDT, 0)
	assert.Nil(t, err)
	t.Log(trades)
}
func TestKuCoinerrGetAccount(t *testing.T) {
	skipKey(t)
	acc, err := kc.GetAccount()
	assert.Nil(t, err)
	t.Log(acc)
}
