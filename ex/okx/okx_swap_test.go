package okx

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestOKExV5Swap_GetFutureTicker(t *testing.T) {
	swap := NewOKExV5Swap(&wstrader.APIConfig{
		HttpClient:    http.DefaultClient,
		ApiKey:        apiKey,
		ApiSecretKey:  apiSecretkey,
		ApiPassphrase: "",
		Lever:         0,
	})
	res, err := swap.GetFutureTicker(cons.BTC_USDT, cons.SWAP_CONTRACT)
	assert.Nil(t, err)
	t.Log(res)
}
func TestOKExV5Swap_GetFutureDepth(t *testing.T) {
	swap := NewOKExV5Swap(&wstrader.APIConfig{
		HttpClient: http.DefaultClient,
	})
	dep, err := swap.GetFutureDepth(cons.BTC_USDT, cons.SWAP_CONTRACT, 2)
	assert.Nil(t, err)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}
func TestOKExV5Swap_GetKlineRecords(t *testing.T) {
	swap := NewOKExV5Swap(&wstrader.APIConfig{
		HttpClient: http.DefaultClient,
	})
	klines, err := swap.GetKlineRecords(cons.SWAP_CONTRACT, cons.BTC_USDT, cons.KLINE_PERIOD_1H, 2)
	assert.Nil(t, err)
	for _, k := range klines {
		t.Logf("%+v", k.Kline)
	}
}
