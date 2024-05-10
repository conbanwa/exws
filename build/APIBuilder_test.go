package build

import (
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"log"
	"testing"
	"time"

	"github.com/conbanwa/logs"
	"github.com/stretchr/testify/assert"
)

var builder = NewAPIBuilder()

func init() {
	logs.Log.Level = logs.L_DEBUG
}
func TestAPIBuilder_Build(t *testing.T) {
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.OKCOIN_COM).String(), cons.OKCOIN_COM)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.HUOBI_PRO).String(), cons.HUOBI_PRO)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.ZB).String(), cons.ZB)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.BIGONE).String(), cons.BIGONE)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.OKEX).String(), cons.OKEX)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.POLONIEX).String(), cons.POLONIEX)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.KRAKEN).String(), cons.KRAKEN)
	assert.Equal(t, builder.APIKey("").APISecretkey("").Build(cons.FCOIN_MARGIN).String(), cons.FCOIN_MARGIN)
	assert.Equal(t, builder.APIKey("").APISecretkey("").BuildFuture(cons.HBDM).String(), cons.HBDM)
}
func TestAPIBuilder_BuildSpotWs(t *testing.T) {
	//os.Setenv("HTTPS_PROXY" , "socks5://"+config.PROXY)
	wsApi, _ := builder.BuildSpotWs(cons.OKEX_V3)
	wsApi.DepthCallback(func(depth *wstrader.Depth) {
		log.Println(depth)
	})
	wsApi.SubscribeDepth(cons.BTC_USDT)
	time.Sleep(time.Minute)
}
func TestAPIBuilder_BuildFuturesWs(t *testing.T) {
	//os.Setenv("HTTPS_PROXY" , "socks5://"+config.PROXY)
	wsApi, _ := builder.BuildFuturesWs(cons.OKEX_V3)
	wsApi.DepthCallback(func(depth *wstrader.Depth) {
		log.Println(depth)
	})
	wsApi.SubscribeDepth(cons.BTC_USD, cons.QUARTER_CONTRACT)
	time.Sleep(time.Minute)
}
