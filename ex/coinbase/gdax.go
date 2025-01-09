package coinbase

import (
	"errors"
	"fmt"
	"github.com/conbanwa/num"
	"github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/web"
	"net/http"
	"sort"

	"github.com/conbanwa/logs"
)

type Coinbase struct {
	httpClient *http.Client
	baseUrl,
	accessKey,
	secretKey string
}

func New(client *http.Client, accesskey, secretkey string) *Coinbase {
	return &Coinbase{client, "https://api.exchange.coinbase.com", accesskey, secretkey}
}
func (c *Coinbase) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	panic("not implement")
}
func (c *Coinbase) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	panic("not implement")
}
func (c *Coinbase) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (c *Coinbase) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (c *Coinbase) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implement")
}
func (c *Coinbase) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (c *Coinbase) GetUnfinishedOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}
func (c *Coinbase) GetOrderHistorys(currency CurrencyPair, optional ...OptionalParameter) ([]Order, error) {
	panic("not implement")
}
func (c *Coinbase) GetAccount() (*Account, error) {
	panic("not implement")
}
func (c *Coinbase) GetTicker(currency CurrencyPair) (*Ticker, error) {
	resp, err := web.HttpGet(c.httpClient, fmt.Sprintf("%s/products/%s/ticker", c.baseUrl, currency.ToSymbol("-")))
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}
	return &Ticker{
		Last: num.ToFloat64(resp["price"]),
		Sell: num.ToFloat64(resp["ask"]),
		Buy:  num.ToFloat64(resp["bid"]),
		Vol:  num.ToFloat64(resp["volume"]),
	}, nil
}
func (c *Coinbase) Get24HStats(pair CurrencyPair) (*Ticker, error) {
	resp, err := web.HttpGet(c.httpClient, fmt.Sprintf("%s/products/%s/stats", c.baseUrl, pair.ToSymbol("-")))
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}
	return &Ticker{
		High: num.ToFloat64(resp["high"]),
		Low:  num.ToFloat64(resp["low"]),
		Vol:  num.ToFloat64(resp["volume"]),
		Last: num.ToFloat64(resp["last"]),
	}, nil
}
func (c *Coinbase) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	var level = 2
	if size == 1 {
		level = 1
	}
	resp, err := web.HttpGet(c.httpClient, fmt.Sprintf("%s/products/%s/book?level=%d", c.baseUrl, currency.ToSymbol("-"), level))
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}
	bids, _ := resp["bids"].([]any)
	asks, _ := resp["asks"].([]any)
	dep := new(Depth)
	for _, v := range bids {
		r := v.([]any)
		dep.BidList = append(dep.BidList, DepthRecord{Price: num.ToFloat64(r[0]), Amount: num.ToFloat64(r[1])})
	}
	for _, v := range asks {
		r := v.([]any)
		dep.AskList = append(dep.AskList, DepthRecord{Price: num.ToFloat64(r[0]), Amount: num.ToFloat64(r[1])})
	}
	sort.Sort(sort.Reverse(dep.AskList))
	return dep, nil
}
func (c *Coinbase) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, opt ...OptionalParameter) ([]Kline, error) {
	urlpath := fmt.Sprintf("%s/products/%s/candles", c.baseUrl, currency.AdaptUsdtToUsd().ToSymbol("-"))
	granularity := -1
	switch period {
	case KLINE_PERIOD_1MIN:
		granularity = 60
	case KLINE_PERIOD_5MIN:
		granularity = 300
	case KLINE_PERIOD_15MIN:
		granularity = 900
	case KLINE_PERIOD_1H, KLINE_PERIOD_60MIN:
		granularity = 3600
	case KLINE_PERIOD_6H:
		granularity = 21600
	case KLINE_PERIOD_1DAY:
		granularity = 86400
	default:
		return nil, errors.New("unsupport the kline period")
	}
	urlpath += fmt.Sprintf("?granularity=%d", granularity)
	resp, err := web.HttpGet3(c.httpClient, urlpath, map[string]string{})
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}
	var klines []wstrader.Kline
	for i := 0; i < len(resp); i++ {
		k, is := resp[i].([]any)
		if !is {
			logs.E("data format err data =", resp[i])
			continue
		}
		klines = append(klines, wstrader.Kline{
			Pair:      currency,
			Timestamp: num.ToInt[int64](k[0]),
			Low:       num.ToFloat64(k[1]),
			High:      num.ToFloat64(k[2]),
			Open:      num.ToFloat64(k[3]),
			Close:     num.ToFloat64(k[4]),
			Vol:       num.ToFloat64(k[5]),
		})
	}
	return klines, nil
}

// 非个人，整个交易所的交易记录
func (c *Coinbase) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}
func (c *Coinbase) String() string {
	return COINBASE
}
