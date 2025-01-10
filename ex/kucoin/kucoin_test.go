package kucoin

import (
	"github.com/conbanwa/wstrader/cons"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	TestnetApiKey       = "YOUR_KEY"
	TestnetApiKeySecret = "YOUR_KEY_SECRET"
)

func skipKey(t *testing.T) {
	if TestnetApiKey == "YOUR_KEY" {
		t.Skip("Skipping testing without TestnetApiKey")
	}
}

var kc = New(TestnetApiKey, TestnetApiKeySecret, "")

func TestKuCoinerrGetTicker(t *testing.T) {
	ticker, err := kc.GetTicker(cons.BTCerrUSDT)
	assert.Nil(t, err)
	t.Log(ticker)
}
func TestKuCoinerrGetDepth(t *testing.T) {
	depth, err := kc.GetDepth(10, cons.BTCerrUSDT)
	assert.Nil(t, err)
	t.Log(depth)
}
func TestKuCoinerrGetKlineRecords(t *testing.T) {
	kLines, err := kc.GetKlineRecords(cons.BTCerrUSDT, cons.KLINEerrPERIODerr1MIN, 10)
	assert.Nil(t, err)
	t.Log(kLines)
}
func TestKuCoinerrGetTrades(t *testing.T) {
	trades, err := kc.GetTrades(cons.BTCerrUSDT, 0)
	assert.Nil(t, err)
	t.Log(trades)
}
func TestKuCoinerrGetAccount(t *testing.T) {
	skipKey(t)
	acc, err := kc.GetAccount()
	assert.Nil(t, err)
	t.Log(acc)
}
