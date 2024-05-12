package bittrex

import (
	"errors"
	"fmt"
	"net/http"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	. "github.com/conbanwa/wstrader/web"
	"sort"

	"github.com/conbanwa/num"
)

type Bittrex struct {
	client *http.Client
	baseUrl,
	accesskey,
	secretkey string
}

func New(client *http.Client, accesskey, secretkey string) *Bittrex {
	return &Bittrex{client: client, accesskey: accesskey, secretkey: secretkey, baseUrl: "https://bittrex.com/api/v1.1"}
}
func (Bittrex *Bittrex) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	panic("not implement")
}
func (Bittrex *Bittrex) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	panic("not implement")
}
func (Bittrex *Bittrex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (Bittrex *Bittrex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (Bittrex *Bittrex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implement")
}
func (Bittrex *Bittrex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (Bittrex *Bittrex) GetUnfinishedOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implement")
}
func (Bittrex *Bittrex) GetOrderHistorys(currency CurrencyPair, optional ...OptionalParameter) ([]Order, error) {
	panic("not implement")
}
func (Bittrex *Bittrex) GetAccount() (*Account, error) {
	panic("not implement")
}
func (Bittrex *Bittrex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	resp, err := HttpGet(Bittrex.client, fmt.Sprintf("%s/public/getmarketsummary?market=%s", Bittrex.baseUrl, currency.ToSymbol2("-")))
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}
	result, _ := resp["result"].([]any)
	if len(result) <= 0 {
		return nil, API_ERR
	}
	tickermap := result[0].(map[string]any)
	return &Ticker{
		Last: num.ToFloat64(tickermap["Last"]),
		Sell: num.ToFloat64(tickermap["Ask"]),
		Buy:  num.ToFloat64(tickermap["Bid"]),
		Low:  num.ToFloat64(tickermap["Low"]),
		High: num.ToFloat64(tickermap["High"]),
		Vol:  num.ToFloat64(tickermap["Volume"]),
	}, nil
}
func (Bittrex *Bittrex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	resp, err := HttpGet(Bittrex.client, fmt.Sprintf("%s/public/getorderbook?market=%s&type=both", Bittrex.baseUrl, currency.ToSymbol2("-")))
	if err != nil {
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return nil, errCode
	}
	result, err2 := resp["result"].(map[string]any)
	if err2 != true {
		return nil, errors.New(resp["message"].(string))
	}
	bids, _ := result["buy"].([]any)
	asks, _ := result["sell"].([]any)
	dep := new(Depth)
	for _, v := range bids {
		r := v.(map[string]any)
		dep.BidList = append(dep.BidList, DepthRecord{Price: num.ToFloat64(r["Rate"]), Amount: num.ToFloat64(r["Quantity"])})
	}
	for _, v := range asks {
		r := v.(map[string]any)
		dep.AskList = append(dep.AskList, DepthRecord{Price: num.ToFloat64(r["Rate"]), Amount: num.ToFloat64(r["Quantity"])})
	}
	sort.Sort(sort.Reverse(dep.AskList))
	return dep, nil
}
func (Bittrex *Bittrex) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, opt ...OptionalParameter) ([]Kline, error) {
	panic("not implement")
}

// 非个人，整个交易所的交易记录
func (Bittrex *Bittrex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}
func (Bittrex *Bittrex) String() string {
	return BITTREX
}
