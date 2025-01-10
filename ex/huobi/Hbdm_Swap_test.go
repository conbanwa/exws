package huobi

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

var swap *HbdmSwap

func init() {
	swap = NewHbdmSwap(&wstrader.APIConfig{
		HttpClient:   http.DefaultClient,
		Endpoint:     "https://api.btcgateway.pro",
		ApiKey:       TestKey,
		ApiSecretKey: TestSecret,
		Lever:        5,
	})
}
func TestHbdmSwap_GetFutureTicker(t *testing.T) {
	res, err := swap.GetFutureTicker(cons.BTC_USD, cons.SWAP_CONTRACT)
	assert.Nil(t, err)
	t.Log(res)
}
func TestHbdmSwap_GetFutureDepth(t *testing.T) {
	dep, err := swap.GetFutureDepth(cons.BTC_USD, cons.SWAP_CONTRACT, 5)
	assert.Nil(t, err)
	t.Log(dep.AskList)
	t.Log(dep.BidList)
}
func TestHbdmSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(swap.GetFutureUserinfo(cons.NewCurrencyPair2("DOT_USD")))
}
func TestHbdmSwap_GetFuturePosition(t *testing.T) {
	t.Log(swap.GetFuturePosition(cons.NewCurrencyPair2("DOT_USD"), cons.SWAP_CONTRACT))
}
func TestHbdmSwap_LimitFuturesOrder(t *testing.T) {
	//784115347040780289
	t.Log(swap.LimitFuturesOrder(cons.NewCurrencyPair2("DOT_USD"), cons.SWAP_CONTRACT, "6.5", "1", cons.OPEN_SELL))
}
func TestHbdmSwap_FutureCancelOrder(t *testing.T) {
	t.Log(swap.FutureCancelOrder(cons.NewCurrencyPair2("DOT_USD"), cons.SWAP_CONTRACT, "784118017750929408"))
}
func TestHbdmSwap_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(swap.GetUnfinishFutureOrders(cons.NewCurrencyPair2("DOT_USD"), cons.SWAP_CONTRACT))
}
func TestHbdmSwap_GetFutureOrder(t *testing.T) {
	t.Log(swap.GetFutureOrder("784118017750929408", cons.NewCurrencyPair2("DOT_USD"), cons.SWAP_CONTRACT))
}
func TestHbdmSwap_GetFutureOrderHistory(t *testing.T) {
	t.Log(swap.GetFutureOrderHistory(cons.NewCurrencyPair2("KSM_USD"), cons.SWAP_CONTRACT,
		wstrader.OptionalParameter{}.Optional("start_time", time.Now().Add(-5*24*time.Hour).Unix()*1000),
		wstrader.OptionalParameter{}.Optional("end_time", time.Now().Unix()*1000)))
}
