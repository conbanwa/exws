package gateio

import (
	"fmt"
	"github.com/conbanwa/num"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	. "github.com/conbanwa/wstrader/web"
	"net/http"
	"sort"
	"strings"
)

var marketBaseUrl = "http://data.gateapi.io/api2/1"

type Gate struct {
	client *http.Client
	accesskey,
	secretkey, phrase string
}

func New(client *http.Client, accesskey, secretkey string) *Gate {
	return &Gate{client: client, accesskey: accesskey, secretkey: secretkey}
}
func (g *Gate) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	panic("not implement")
}
func (g *Gate) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	panic("not implement")
}
func (g *Gate) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Gate) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Gate) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implement")
}
func (g *Gate) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (g *Gate) GetUnfinishedOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}
func (g *Gate) GetOrderHistorys(currency CurrencyPair, para ...OptionalParameter) ([]Order, error) {
	panic("not implement")
}
func (g *Gate) GetAccount() (*Account, error) {
	// panic("not implement")
	return nil, nil
}
func (g *Gate) GetTicker(currency CurrencyPair) (*Ticker, error) {
	uri := fmt.Sprintf("%s/ticker/%s", marketBaseUrl, strings.ToLower(currency.ToSymbol("_")))
	resp, err := HttpGet(g.client, uri)
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}
	return &Ticker{
		Last: num.ToFloat64(resp["last"]),
		Sell: num.ToFloat64(resp["lowestAsk"]),
		Buy:  num.ToFloat64(resp["highestBid"]),
		High: num.ToFloat64(resp["high24hr"]),
		Low:  num.ToFloat64(resp["low24hr"]),
		Vol:  num.ToFloat64(resp["quoteVolume"]),
	}, nil
}
func (g *Gate) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	resp, err := HttpGet(g.client, fmt.Sprintf("%s/orderBook/%s", marketBaseUrl, currency.ToSymbol("_")))
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
func (g *Gate) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]Kline, error) {
	panic("not implement")
}

// 非个人，整个交易所的交易记录
func (g *Gate) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}
func (g *Gate) String() string {
	return GATEIO
}
