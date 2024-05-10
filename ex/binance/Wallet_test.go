package binance

import (
	"net/http"
	"qa3/wstrader"
	"qa3/wstrader/cons"
	"testing"
)

var wallet *Wallet

func init() {
	wallet = NewWallet(&wstrader.APIConfig{
		HttpClient:   http.DefaultClient,
		ApiKey:       "",
		ApiSecretKey: "",
	})
}
func TestWallet_Transfer(t *testing.T) {
	t.Log(wallet.Transfer(wstrader.TransferParameter{
		Currency: "USDT",
		From:     cons.SPOT,
		To:       cons.SWAP_USDT,
		Amount:   100,
	}))
}
