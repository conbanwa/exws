package okx

import (
	"fmt"
	"github.com/conbanwa/num"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/util"
	"math"
	"net/url"
	"time"
)

type V5Spot struct {
	*OKX
}

func NewOKExV5Spot(config *APIConfig) *V5Spot {
	if config.Endpoint == "" {
		config.Endpoint = v5RestBaseUrl
	}
	okx := &V5Spot{OKX: NewOKExV5(config)}
	return okx
}

// private API
func (ok *V5Spot) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	ty := "limit"
	if len(opt) > 0 {
		ty = opt[0].String()
	}
	response, err := ok.CreateOrder(&CreateOrderParam{
		Symbol:    currency.ToSymbol("-"),
		TradeMode: "cash",
		Side:      "buy",
		OrderType: ty,
		Size:      amount,
		Price:     price,
	})
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		Price:    num.ToFloat64(price),
		Amount:   num.ToFloat64(amount),
		Cid:      response.ClientOrdId,
		OrderID2: response.OrdId,
	}, nil
}
func (ok *V5Spot) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	ty := "limit"
	if len(opt) > 0 {
		ty = opt[0].String()
	}
	response, err := ok.CreateOrder(&CreateOrderParam{
		Symbol:    currency.ToSymbol("-"),
		TradeMode: "cash",
		Side:      "sell",
		OrderType: ty,
		Size:      amount,
		Price:     price,
	})
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		Price:    num.ToFloat64(price),
		Amount:   num.ToFloat64(amount),
		Cid:      response.ClientOrdId,
		OrderID2: response.OrdId,
	}, nil
}
func (ok *V5Spot) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	response, err := ok.CreateOrder(&CreateOrderParam{
		Symbol:    currency.ToSymbol("-"),
		TradeMode: "cash",
		Side:      "buy",
		OrderType: "market",
		Size:      amount,
	})
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		Amount:   num.ToFloat64(amount),
		Cid:      response.ClientOrdId,
		OrderID2: response.OrdId,
	}, nil
}
func (ok *V5Spot) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	response, err := ok.CreateOrder(&CreateOrderParam{
		Symbol:    currency.ToSymbol("-"),
		TradeMode: "cash",
		Side:      "sell",
		OrderType: "market",
		Size:      amount,
	})
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		Amount:   num.ToFloat64(amount),
		Cid:      response.ClientOrdId,
		OrderID2: response.OrdId,
	}, nil
}
func (ok *V5Spot) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	_, err := ok.CancelOrderV5(currency.ToSymbol("-"), orderId, "")
	if err != nil {
		return false, err
	}
	return true, nil
}
func (ok *V5Spot) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	response, err := ok.GetOrderV5(currency.ToSymbol("-"), orderId, "")
	if err != nil {
		return nil, err
	}
	status := ORDER_UNFINISH
	switch response.State {
	case "canceled":
		status = ORDER_CANCEL
	case "live":
		status = ORDER_UNFINISH
	case "partially_filled":
		status = ORDER_PART_FINISH
	case "filled":
		status = ORDER_FINISH
	default:
		status = ORDER_UNFINISH
	}
	side := BUY
	if response.Side == "sell" || response.Side == "SELL" {
		side = SELL
	}
	return &Order{
		Price:        response.Px,
		Amount:       response.Sz,
		AvgPrice:     num.ToFloat64(response.AvgPx),
		DealAmount:   num.ToFloat64(response.AccFillSz),
		Fee:          response.Fee,
		Cid:          response.ClOrdID,
		OrderID2:     response.OrdID,
		Status:       status,
		Currency:     currency,
		Side:         side,
		Type:         response.OrdType,
		OrderTime:    response.CTime,
		FinishedTime: response.UTime,
	}, nil
}
func (ok *V5Spot) GetUnfinishedOrders(currency CurrencyPair) ([]Order, error) {
	response, err := ok.GetPendingOrders(&PendingOrderParam{
		InstType: "SPOT",
		InstId:   currency.ToSymbol("-"),
	})
	if err != nil {
		return nil, err
	}
	orders := make([]Order, 0)
	for _, v := range response {
		status := ORDER_UNFINISH
		switch v.State {
		case "canceled":
			status = ORDER_CANCEL
		case "live":
			status = ORDER_UNFINISH
		case "partially_filled":
			status = ORDER_PART_FINISH
		case "filled":
			status = ORDER_FINISH
		default:
			status = ORDER_UNFINISH
		}
		side := BUY
		if v.Side == "sell" || v.Side == "SELL" {
			side = SELL
		}
		orders = append(orders, Order{
			Price:        v.Px,
			Amount:       v.Sz,
			AvgPrice:     num.ToFloat64(v.AvgPx),
			DealAmount:   num.ToFloat64(v.AccFillSz),
			Fee:          v.Fee,
			Cid:          v.ClOrdID,
			OrderID2:     v.OrdID,
			Status:       status,
			Currency:     currency,
			Side:         side,
			Type:         v.OrdType,
			OrderTime:    v.CTime,
			FinishedTime: v.UTime,
		})
	}
	return orders, nil
}
func (ok *V5Spot) GetOrderHistorys(currency CurrencyPair, opt ...OptionalParameter) ([]Order, error) {
	response, err := ok.GetOrderHistory(
		"SPOT",
		"", //currency.ToSymbol("-"),
		"", "", "", "",
	)
	if err != nil {
		return nil, err
	}
	orders := make([]Order, 0)
	for _, v := range response {
		status := ORDER_UNFINISH
		switch v.State {
		case "canceled":
			status = ORDER_CANCEL
		case "live":
			status = ORDER_UNFINISH
		case "partially_filled":
			status = ORDER_PART_FINISH
		case "filled":
			status = ORDER_FINISH
		default:
			status = ORDER_UNFINISH
		}
		side := BUY
		if v.Side == "sell" || v.Side == "SELL" {
			side = SELL
		}
		orders = append(orders, Order{
			Price:        v.Px,
			Amount:       v.Sz,
			AvgPrice:     num.ToFloat64(v.AvgPx),
			DealAmount:   num.ToFloat64(v.AccFillSz),
			Fee:          v.Fee,
			Cid:          v.ClOrdID,
			OrderID2:     v.OrdID,
			Status:       status,
			Currency:     NewCurrencyPair3(v.InstID, "-"),
			Side:         side,
			Type:         v.OrdType,
			OrderTime:    v.CTime,
			FinishedTime: v.UTime,
		})
	}
	return orders, nil
}
func (ok *V5Spot) GetAccount() (*Account, error) {
	response, err := ok.GetAccountBalances("")
	if err != nil {
		return nil, err
	}
	account := &Account{
		SubAccounts: make(map[Currency]SubAccount, 2)}
	for _, itm := range response.Details {
		currency := NewCurrency(itm.Currency, "")
		account.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			ForzenAmount: num.ToFloat64(itm.Frozen),
			Amount:       math.Max(num.ToFloat64(itm.Available), num.ToFloat64(itm.AvailEq)),
		}
	}
	return account, nil
}

// public API
func (ok *V5Spot) GetTicker(currency CurrencyPair) (*Ticker, error) {
	ticker, err := ok.GetTickerV5(currency.ToSymbol("-"))
	if err != nil {
		return nil, err
	}
	return &Ticker{
		Pair: currency,
		Last: ticker.Last,
		Buy:  ticker.BuyPrice,
		Sell: ticker.SellPrice,
		High: ticker.High,
		Low:  ticker.Low,
		Vol:  ticker.Vol,
		Date: ticker.Timestamp,
	}, nil
}
func (ok *V5Spot) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	d, err := ok.GetDepthV5(currency.ToSymbol("-"), size)
	if err != nil {
		return nil, err
	}
	depth := &Depth{}
	for _, ask := range d.Asks {
		depth.AskList = append(depth.AskList, DepthRecord{Price: num.ToFloat64(ask[0]), Amount: num.ToFloat64(ask[1])})
	}
	for _, bid := range d.Bids {
		depth.BidList = append(depth.BidList, DepthRecord{Price: num.ToFloat64(bid[0]), Amount: num.ToFloat64(bid[1])})
	}
	depth.UTime = time.Unix(0, int64(d.Timestamp)*1000000)
	depth.Pair = currency
	return depth, nil
}
func (ok *V5Spot) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]Kline, error) {
	// [1m/3m/5m/15m/30m/1H/2H/4H/6H/12H/1D/1W/1M/3M/6M/1Y]
	param := &url.Values{}
	param.Set("limit", fmt.Sprint(size))
	util.MergeOptionalParameter(param, optional...)
	kl, err := ok.GetKlineRecordsV5(currency.ToSymbol("-"), period, param)
	if err != nil {
		return nil, err
	}
	klines := make([]Kline, 0)
	for _, k := range kl {
		klines = append(klines, Kline{
			Pair:      currency,
			Timestamp: num.ToInt[int64](k[0]),
			Open:      num.ToFloat64(k[1]),
			High:      num.ToFloat64(k[2]),
			Low:       num.ToFloat64(k[3]),
			Close:     num.ToFloat64(k[4]),
			Vol:       num.ToFloat64(k[5]),
		})
	}
	return klines, nil
}

// 非个人，整个交易所的交易记录
func (ok *V5Spot) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not support")
}
func (ok *V5Spot) String() string {
	return ok.ExchangeName()
}
