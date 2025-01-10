package binance

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"net/http"
	"testing"
)

var wallet *Wallet

func init() {
	wallet = NewWallet(&wstrader.APIConfig{
		HttpClient:   http.DefaultClient,
		ApiKey:       TestnetApiKey,
		ApiSecretKey: TestnetApiKeySecret,
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
