package binance

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/cons"
	"net/http"
	"testing"
)

var wallet *Wallet

func init() {
	wallet = NewWallet(&exws.APIConfig{
		HttpClient:   http.DefaultClient,
		ApiKey:       apiKey,
		ApiSecretKey: apiSecretkey,
	})
}
func TestWallet_Transfer(t *testing.T) {
	t.Log(wallet.Transfer(exws.TransferParameter{
		Currency: "USDT",
		From:     cons.SPOT,
		To:       cons.SWAP_USDT,
		Amount:   100,
	}))
}
