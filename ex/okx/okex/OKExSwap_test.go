package okex

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/config"
	"github.com/conbanwa/wstrader/cons"
	"net/http"
	"net/url"
	"testing"
	"time"
)

var configs = &wstrader.APIConfig{
	HttpClient: &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return &url.URL{
					Scheme: "socks5",
					Host:   config.Proxy}, nil
			},
		},
	},
	Endpoint:      "https://www.okx.com",
	ApiKey:        "",
	ApiSecretKey:  "",
	ApiPassphrase: "",
}
var okExSwap = NewOKExSwap(configs)

func TestOKExSwap_GetFutureUserinfo(t *testing.T) {
	t.Log(okExSwap.GetFutureUserinfo())
}
func TestOKExSwap_PlaceFutureOrder(t *testing.T) {
	t.Log(okExSwap.PlaceFutureOrder(cons.BTC_USDT, cons.SWAP_CONTRACT, "10000", "1", cons.OPEN_BUY, 0, 0))
}
func TestOKExSwap_PlaceFutureOrder2(t *testing.T) {
	t.Log(okExSwap.PlaceFutureOrder2(cons.BTC_USDT, cons.SWAP_CONTRACT, "10000", "1", cons.OPEN_BUY, 0, cons.Ioc))
}
func TestOKExSwap_FutureCancelOrder(t *testing.T) {
	t.Log(okExSwap.FutureCancelOrder(cons.BTC_USDT, cons.SWAP_CONTRACT, "309935122485305344"))
}
func TestOKExSwap_GetFutureOrder(t *testing.T) {
	t.Log(okExSwap.GetFutureOrder("581084124456583168", cons.BTC_USDT, cons.SWAP_CONTRACT))
}
func TestOKExSwap_GetFuturePosition(t *testing.T) {
	t.Log(okExSwap.GetFuturePosition(cons.BTC_USD, cons.SWAP_CONTRACT))
}
func TestOKExSwap_GetFutureDepth(t *testing.T) {
	t.Log(okExSwap.GetFutureDepth(cons.LTC_USD, cons.SWAP_CONTRACT, 10))
}
func TestOKExSwap_GetFutureTicker(t *testing.T) {
	t.Log(okExSwap.GetFutureTicker(cons.BTC_USD, cons.SWAP_CONTRACT))
}
func TestOKExSwap_GetUnfinishFutureOrders(t *testing.T) {
	ords, _ := okExSwap.GetUnfinishFutureOrders(cons.XRP_USD, cons.SWAP_CONTRACT)
	for _, ord := range ords {
		t.Log(ord.OrderID2, ord.ClientOid)
	}
}
func TestOKExSwap_GetHistoricalFunding(t *testing.T) {
	for i := 1; ; i++ {
		funding, err := okExSwap.GetHistoricalFunding(cons.SWAP_CONTRACT, cons.BTC_USD, i)
		t.Log(err, len(funding))
	}
}
func TestOKExSwap_GetKlineRecords(t *testing.T) {
	since := time.Now().Add(-24 * time.Hour).Unix()
	kline, err := okExSwap.GetKlineRecords(cons.SWAP_CONTRACT, cons.BTC_USD, cons.KLINE_PERIOD_4H, 0, wstrader.OptionalParameter{"since": since})
	t.Log(err, kline[0].Kline)
}
func TestOKExSwap_GetKlineRecords2(t *testing.T) {
	start := time.Now().Add(time.Minute * -30).UTC().Format(time.RFC3339)
	t.Log(start)
	kline, err := okExSwap.GetKlineRecords2(cons.SWAP_CONTRACT, cons.BTC_USDT, start, "", "900")
	t.Log(err, kline[0].Kline)
}
func TestOKExSwap_GetInstruments(t *testing.T) {
	t.Log(okExSwap.GetInstruments())
}
func TestOKExSwap_SetMarginLevel(t *testing.T) {
	t.Log(okExSwap.SetMarginLevel(cons.EOS_USDT, 5, 3))
}
func TestOKExSwap_GetMarginLevel(t *testing.T) {
	t.Log(okExSwap.GetMarginLevel(cons.EOS_USDT))
}
func TestOKExSwap_GetFutureAccountInfo(t *testing.T) {
	t.Log(okExSwap.GetFutureAccountInfo(cons.BTC_USDT))
}
func TestOKExSwap_PlaceFutureAlgoOrder(t *testing.T) {
	ord := &wstrader.FutureOrder{
		ContractName: cons.SWAP_CONTRACT,
		Currency:     cons.BTC_USD,
		OType:        2, //开空
		OrderType:    1, //1：止盈止损 2：跟踪委托 3：冰山委托 4：时间加权
		Price:        9877,
		Amount:       1,
		TriggerPrice: 9877,
		AlgoType:     1,
	}
	t.Log(okExSwap.PlaceFutureAlgoOrder(ord))
}
func TestOKExSwap_FutureCancelAlgoOrder(t *testing.T) {
	t.Log(okExSwap.FutureCancelAlgoOrder(cons.BTC_USD, []string{"309935122485305344"}))
}
func TestOKExSwap_GetFutureAlgoOrders(t *testing.T) {
	t.Log(okExSwap.GetFutureAlgoOrders("", "2", cons.BTC_USD))
}
