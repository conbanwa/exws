package okx

import (
	"github.com/conbanwa/exws"
	"github.com/conbanwa/exws/config"
	"github.com/conbanwa/exws/cons"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func newOKExV5Client() *OKX {
	return NewOKExV5(&exws.APIConfig{
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
		ApiKey:        apiKey,
		ApiSecretKey:  apiSecretkey,
		ApiPassphrase: "",
	})
}
func TestOKExV5_GetTicker(t *testing.T) {
	o := newOKExV5Client()
	res, err := o.GetTickerV5("BTC-USD-SWAP")
	assert.Nil(t, err)
	t.Log(res)
}
func TestOKExV5_GetDepth(t *testing.T) {
	o := newOKExV5Client()
	res, err := o.GetDepthV5("BTC-USD-SWAP", 0)
	assert.Nil(t, err)
	t.Log(res)
}
func TestOKExV5_GetKlineRecordsV5(t *testing.T) {
	o := newOKExV5Client()
	res, err := o.GetKlineRecordsV5("BTC-USD-SWAP", cons.KLINE_PERIOD_5MIN, &url.Values{})
	assert.Nil(t, err)
	t.Log("total: ", len(res), res[0], res[1], res[2])
}
