package okx

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/config"
	"github.com/conbanwa/wstrader/cons"
	"net/http"
	"net/url"
	"testing"
)

const (
	TestnetApiKey       = "YOUR_KEY"
	TestnetApiKeySecret = "YOUR_KEY_SECRET"
)

var c = newOKExV5SpotClient()

func skipKey(t *testing.T) {
	if TestnetApiKey == "YOUR_KEY" {
		t.Skip("Skipping testing without TestnetApiKey")
	}
}

func newOKExV5SpotClient() *V5Spot {
	return NewOKExV5Spot(&wstrader.APIConfig{
		HttpClient: &http.Client{
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					return &url.URL{
						Scheme: "socks5",
						Host:   config.Proxy}, nil
				},
			},
		},
		Endpoint:      baseUrl,
		ApiKey:        TestnetApiKey,
		ApiSecretKey:  TestnetApiKeySecret,
		ApiPassphrase: "",
	}
	assert.Nil(t, err)
	t.Log(res)
}
func init() {
}
func TestOKExV5Spot_GetTicker(t *testing.T) {
	res, err := c.GetTicker(cons.BTC_USDT)
	assert.Nil(t, err)
	t.Log(res)
}
func TestOKExV5Spot_GetDepth(t *testing.T) {
	res, err := c.GetDepth(5, cons.BTC_USDT)
	assert.Nil(t, err)
	t.Log(res)
}
func TestOKExV5SpotGetKlineRecords(t *testing.T) {
	res, err := c.GetKlineRecords(cons.BTC_USDT, cons.KLINE_PERIOD_1MIN, 10)
	assert.Nil(t, err)
	t.Log(res)
}
func TestOKExV5Spot_LimitBuy(t *testing.T) {
	skipKey(t)
	res, err := c.LimitBuy("1", "1.0", cons.XRP_USDT))
	//{"code":"0","data":[{"clOrdId":"0bf60374efe445BC258eddf46df044c3","ordId":"305267682086109184","sCode":"0","sMsg":"","tag":""}],"msg":""}}
	assert.Nil(t, err)
	t.Log(res)
}
func TestOKExV5Spot_CancelOrder(t *testing.T) {
	skipKey(t)
	res, err := c.CancelOrder("305267682086109184", cons.XRP_USDT)
	assert.Nil(t, err)
	t.Log(res)
}
func TestOKExV5Spot_GetUnfinishOrders(t *testing.T) {
	skipKey(t)
	res, err := c.GetUnfinishedOrders(cons.XRP_USDT)
	assert.Nil(t, err)
	t.Log(res)
}
func TestOKExV5Spot_GetOneOrder(t *testing.T) {
	skipKey(t)
	res, err := c.GetOneOrder("305267682086109184", cons.XRP_USDT)
	assert.Nil(t, err)
	t.Log(res)
}
func TestOKExV5Spot_GetAccount(t *testing.T) {
	skipKey(t)
	res, err := c.GetAccount()
	assert.Nil(t, err)
	t.Log(res)
}
