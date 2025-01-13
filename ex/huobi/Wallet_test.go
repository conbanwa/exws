package huobi

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/cons"
	"testing"
)

var wallet *Wallet

func init() {
	wallet = NewWallet(&exws.APIConfig{
		HttpClient:   httpProxyClient,
		ApiKey:       apiKey,
		ApiSecretKey: apiSecretkey,
	})
}
func TestWallet_Transfer(t *testing.T) {
	t.Log(wallet.Transfer(exws.TransferParameter{
		Currency: "BTC",
		From:     cons.SWAP_USDT,
		To:       cons.SPOT,
		Amount:   11,
	}))
}
