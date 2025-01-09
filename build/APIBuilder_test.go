package build

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
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
	assert.Equal(t, builder.APIKey("").APISecretkey("").BuildFuture(cons.HBDM).String(), cons.HBDM)
}
func TestAPIBuilder_BuildSpotWs(t *testing.T) {
	buildSpotWs(t, cons.BINANCE)
	buildSpotWs(t, cons.GATEIO)
	buildSpotWs(t, cons.HUOBI_PRO)
	time.Sleep(time.Minute)
}
func buildSpotWs(t *testing.T, ex string) {
	wsApi, err := builder.BuildSpotWs(ex)
	assert.Nil(t, err)
	wsApi.BBOCallback(func(bbo *q.Bbo) {
		t.Log(bbo)
	})

	wsApi.SubscribeBBO([]string{})
	wsApi.DepthCallback(func(depth *wstrader.Depth) {
		t.Log(depth)
	})
	wsApi.SubscribeDepth(cons.BTC_USDT)
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
