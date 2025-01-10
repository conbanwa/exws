package huobi

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"testing"
)

var wallet *Wallet

func init() {
	wallet = NewWallet(&wstrader.APIConfig{
		HttpClient:   httpProxyClient,
		ApiKey:       apiKey,
		ApiSecretKey: apiSecretkey,
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
