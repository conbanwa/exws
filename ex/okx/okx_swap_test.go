package okx

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"net/http"
	"testing"

	"github.com/conbanwa/logs"
)

func init() {
	logs.Log.Level = logs.L_DEBUG
	//os.Setenv("HTTPS_PROXY", "socks5://" + conf.PROXY) //local socks5 proxy
}
func TestOKExV5Swap_GetFutureTicker(t *testing.T) {
	swap := NewOKExV5Swap(&wstrader.APIConfig{
		HttpClient:    http.DefaultClient,
		ApiKey:        "",
		ApiSecretKey:  "",
		ApiPassphrase: "",
		Lever:         0,
	})
	t.Log(swap.GetFutureTicker(cons.BTC_USDT, cons.SWAP_CONTRACT))
}
func TestOKExV5Swap_GetFutureDepth(t *testing.T) {
	swap := NewOKExV5Swap(&wstrader.APIConfig{
		HttpClient: http.DefaultClient,
	})
	dep, err := swap.GetFutureDepth(cons.BTC_USDT, cons.SWAP_CONTRACT, 2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}
func TestOKExV5Swap_GetKlineRecords(t *testing.T) {
	swap := NewOKExV5Swap(&wstrader.APIConfig{
		HttpClient: http.DefaultClient,
	})
	klines, err := swap.GetKlineRecords(cons.SWAP_CONTRACT, cons.BTC_USDT, cons.KLINE_PERIOD_1H, 2)
	if err != nil {
		t.Error(err)
		return
	}
	for _, k := range klines {
		t.Logf("%+v", k.Kline)
	}
}
