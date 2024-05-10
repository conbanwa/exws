package huobi

import (
	"qa3/wstrader"
	"qa3/wstrader/cons"
	"testing"
)

var wallet *Wallet

func init() {
	wallet = NewWallet(&wstrader.APIConfig{
		HttpClient:   httpProxyClient,
		ApiKey:       "",
		ApiSecretKey: "",
	})
}
func TestWallet_Transfer(t *testing.T) {
	t.Log(wallet.Transfer(wstrader.TransferParameter{
		Currency: "BTC",
		From:     cons.SWAP_USDT,
		To:       cons.SPOT,
		Amount:   11,
	}))
}
