package huobi

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"testing"
	"time"
)

var dm = NewHbdm(&wstrader.APIConfig{
	Endpoint:     "https://api.hbdm.com",
	HttpClient:   httpProxyClient,
	ApiKey:       "1aa7ddfb-2fcaf72e-b1rkuf4drg-dd9a2",
	ApiSecretKey: "9ca51e4d-dcd8e098-18b345c2-8cea0"})

func TestHbdm_GetFutureUserinfo(t *testing.T) {
	t.Log(dm.GetFutureUserinfo())
}
func TestHbdm_GetFuturePosition(t *testing.T) {
	t.Log(dm.GetFuturePosition(cons.BTC_USD, cons.QUARTER_CONTRACT))
}
func TestHbdm_PlaceFutureOrder(t *testing.T) {
	t.Log(dm.PlaceFutureOrder(cons.BTC_USD, cons.QUARTER_CONTRACT, "3800", "1", cons.OPEN_BUY, 0, 20))
}
func TestHbdm_FutureCancelOrder(t *testing.T) {
	t.Log(dm.FutureCancelOrder(cons.BTC_USD, cons.QUARTER_CONTRACT, "6"))
}
func TestHbdm_GetUnfinishFutureOrders(t *testing.T) {
	t.Log(dm.GetUnfinishFutureOrders(cons.BTC_USD, cons.QUARTER_CONTRACT))
}
func TestHbdm_GetFutureOrders(t *testing.T) {
	t.Log(dm.GetFutureOrders([]string{"6", "5"}, cons.BTC_USD, cons.QUARTER_CONTRACT))
}
func TestHbdm_GetFutureOrder(t *testing.T) {
	t.Log(dm.GetFutureOrder("6", cons.BTC_USD, cons.QUARTER_CONTRACT))
}
func TestHbdm_GetFutureTicker(t *testing.T) {
	t.Log(dm.GetFutureTicker(cons.EOS_USD, cons.QUARTER_CONTRACT))
}
func TestHbdm_GetFutureDepth(t *testing.T) {
	dep, err := dm.GetFutureDepth(cons.BTC_USD, cons.QUARTER_CONTRACT, 0)
	t.Log(err)
	t.Logf("%+v\n%+v", dep.AskList, dep.BidList)
}
func TestHbdm_GetFutureIndex(t *testing.T) {
	t.Log(dm.GetFutureIndex(cons.BTC_USD))
}
func TestHbdm_GetFutureEstimatedPrice(t *testing.T) {
	t.Log(dm.GetFutureEstimatedPrice(cons.BTC_USD))
}
func TestHbdm_GetKlineRecords(t *testing.T) {
	klines, _ := dm.GetKlineRecords(cons.QUARTER_CONTRACT, cons.EOS_USD, cons.KLINE_PERIOD_1MIN, 20, wstrader.OptionalParameter{"test": 0})
	for _, k := range klines {
		tt := time.Unix(k.Timestamp, 0)
		t.Log(k.Pair, tt, k.Open, k.Close, k.High, k.Low, k.Vol, k.Vol2)
	}
}
