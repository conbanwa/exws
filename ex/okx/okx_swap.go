package okx

import (
	"fmt"
	"github.com/conbanwa/logs"
	"github.com/conbanwa/num"
	"net/url"
	. "qa3/wstrader"
	. "qa3/wstrader/cons"
	. "qa3/wstrader/q"
	"qa3/wstrader/util"
	"sort"
	"time"
)

type V5Swap struct {
	*OKX
}

func NewOKExV5Swap(config *APIConfig) *V5Swap {
	v5 := new(V5Swap)
	v5.OKX = NewOKExV5(config)
	return v5
}
func (ok *V5Swap) String() string {
	return OKEX_SWAP
}
func (ok *V5Swap) GetFutureEstimatedPrice(currencyPair CurrencyPair) (float64, error) {
	panic("implement me")
}
func (ok *V5Swap) GetFutureTicker(currencyPair CurrencyPair, contractType string) (*Ticker, error) {
	t, err := ok.OKX.GetTickerV5(fmt.Sprintf("%s-SWAP", currencyPair.ToSymbol("-")))
	if err != nil {
		return nil, err
	}
	return &Ticker{
		Pair: currencyPair,
		Last: t.Last,
		Buy:  t.BuyPrice,
		Sell: t.SellPrice,
		High: t.High,
		Low:  t.Low,
		Vol:  t.Vol,
		Date: t.Timestamp,
	}, nil
}
func (ok *V5Swap) GetFutureDepth(currencyPair CurrencyPair, contractType string, size int) (*Depth, error) {
	instId := fmt.Sprintf("%s-SWAP", currencyPair.ToSymbol("-"))
	dep, err := ok.OKX.GetDepthV5(instId, size)
	if err != nil {
		return nil, err
	}
	depth := &Depth{}
	for _, ask := range dep.Asks {
		depth.AskList = append(depth.AskList, DepthRecord{Price: num.ToFloat64(ask[0]), Amount: num.ToFloat64(ask[1])})
	}
	for _, bid := range dep.Bids {
		depth.BidList = append(depth.BidList, DepthRecord{Price: num.ToFloat64(bid[0]), Amount: num.ToFloat64(bid[1])})
	}
	sort.Sort(sort.Reverse(depth.AskList))
	depth.Pair = currencyPair
	depth.UTime = time.Unix(0, int64(dep.Timestamp)*1000000)
	return depth, nil
}
func (ok *V5Swap) GetFutureIndex(currencyPair CurrencyPair) (float64, error) {
	panic("implement me")
}
func (ok *V5Swap) GetFutureUserinfo(currencyPair ...CurrencyPair) (*FutureAccount, error) {
	panic("implement me")
}
func (ok *V5Swap) PlaceFutureOrder(currencyPair CurrencyPair, contractType, price, amount string, openType, matchPrice int, leverRate float64) (string, error) {
	panic("implement me")
}
func (ok *V5Swap) LimitFuturesOrder(currencyPair CurrencyPair, contractType, price, amount string, openType int, opt ...LimitOrderOptionalParameter) (*FutureOrder, error) {
	panic("implement me")
}
func (ok *V5Swap) MarketFuturesOrder(currencyPair CurrencyPair, contractType, amount string, openType int) (*FutureOrder, error) {
	panic("implement me")
}
func (ok *V5Swap) FutureCancelOrder(currencyPair CurrencyPair, contractType, orderId string) (bool, error) {
	panic("implement me")
}
func (ok *V5Swap) GetFuturePosition(currencyPair CurrencyPair, contractType string) ([]FuturePosition, error) {
	panic("implement me")
}
func (ok *V5Swap) GetFutureOrders(orderIds []string, currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("implement me")
}
func (ok *V5Swap) GetFutureOrder(orderId string, currencyPair CurrencyPair, contractType string) (*FutureOrder, error) {
	panic("implement me")
}
func (ok *V5Swap) GetUnfinishFutureOrders(currencyPair CurrencyPair, contractType string) ([]FutureOrder, error) {
	panic("implement me")
}
func (ok *V5Swap) GetFutureOrderHistory(pair CurrencyPair, contractType string, optional ...OptionalParameter) ([]FutureOrder, error) {
	panic("implement me")
}
func (ok *V5Swap) GetFee() (float64, error) {
	panic("implement me")
}
func (ok *V5Swap) GetContractValue(currencyPair CurrencyPair) (float64, error) {
	panic("implement me")
}
func (ok *V5Swap) GetDeliveryTime() (int, int, int, int) {
	panic("implement me")
}
func (ok *V5Swap) GetKlineRecords(contractType string, currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]FutureKline, error) {
	param := &url.Values{}
	param.Set("limit", fmt.Sprint(size))
	util.MergeOptionalParameter(param, optional...)
	data, err := ok.OKX.GetKlineRecordsV5(fmt.Sprintf("%s-SWAP", currency.ToSymbol("-")), period, param)
	if err != nil {
		return nil, err
	}
	logs.D("[okx v5] kline response data: ", data)
	var klines []FutureKline
	for _, item := range data {
		klines = append(klines, FutureKline{
			Kline: &Kline{
				Pair:      currency,
				Timestamp: num.ToInt[int64](item[0]) / 1000,
				Open:      num.ToFloat64(item[1]),
				Close:     num.ToFloat64(item[4]),
				High:      num.ToFloat64(item[2]),
				Low:       num.ToFloat64(item[3]),
				Vol:       num.ToFloat64(item[5]),
			},
		})
	}
	return klines, nil
}
func (ok *V5Swap) GetTrades(contractType string, currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("implement me")
}
