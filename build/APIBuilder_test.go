package build

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var builder = NewAPIBuilder()

func TestAPIBuilder_Build(t *testing.T) {
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.BINANCE).String(), cons.BINANCE)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.BIGONE).String(), cons.BIGONE)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.BITSTAMP).String(), cons.BITSTAMP)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.HUOBI_PRO).String(), cons.HUOBI_PRO)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.OKEX).String(), cons.OKEX)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.POLONIEX).String(), cons.POLONIEX)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.KRAKEN).String(), cons.KRAKEN)
	api, err := builder.APIKey("").APISecretkey("").BuildSpotWs(cons.GATEIO)
	assert.Nil(t, err)
	assert.Equal(t, builder.APIKey("").APISecretkey("").BuildFuture(cons.HBDM).String(), cons.HBDM)
}
func TestAPIBuilder_BuildSpotWs(t *testing.T) {
	wsApi, err := builder.BuildSpotWs(cons.BINANCE)
	assert.Nil(t, err)
	wsApi.DepthCallback(func(depth *wstrader.Depth) {
		t.Log(depth)
	})
	wsApi.SubscribeDepth(cons.BTC_USDT)
	time.Sleep(time.Minute)
}
func TestAPIBuilder_BuildFuturesWs(t *testing.T) {
	wsApi, err := builder.BuildFuturesWs(cons.BINANCE)
	assert.Nil(t, err)
	wsApi.DepthCallback(func(depth *wstrader.Depth) {
		t.Log(depth)
	})
	wsApi.SubscribeDepth(cons.BTC_USD, cons.QUARTER_CONTRACT)
	time.Sleep(time.Minute)
}
