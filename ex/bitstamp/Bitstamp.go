package bitstamp

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/conbanwa/exws"
	. "github.com/conbanwa/exws/cons"
	. "github.com/conbanwa/exws/q"
	. "github.com/conbanwa/exws/util"
	. "github.com/conbanwa/exws/web"
	"github.com/conbanwa/num"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/conbanwa/logs"
)

var (
	BASE_URL = "https://www.bitstamp.net/api/"
)
var _INTERNAL_KLINE_PERIOD_CONVERTER = map[KlinePeriod]string{
	KLINE_PERIOD_1MIN:  "60",
	KLINE_PERIOD_3MIN:  "180",
	KLINE_PERIOD_5MIN:  "300",
	KLINE_PERIOD_15MIN: "900",
	KLINE_PERIOD_30MIN: "1800",
	KLINE_PERIOD_60MIN: "3600",
	KLINE_PERIOD_1H:    "3600",
	KLINE_PERIOD_2H:    "7200",
	KLINE_PERIOD_4H:    "14400",
	KLINE_PERIOD_6H:    "21600",
	// KLINE_PERIOD_8H:     "28800", // Not supported
	KLINE_PERIOD_12H:  "43200",
	KLINE_PERIOD_1DAY: "86400",
	KLINE_PERIOD_3DAY: "259200",
	// KLINE_PERIOD_1WEEK:  "604800", //Not supported
	// KLINE_PERIOD_1MONTH: "1M", // Not supported
}

type Bitstamp struct {
	client *http.Client
	clientId,
	accessKey,
	secretkey string
}

func NewBitstamp(client *http.Client, accessKey, secertkey, clientId string) *Bitstamp {
	return &Bitstamp{client: client, accessKey: accessKey, secretkey: secertkey, clientId: clientId}
}
func (Bitstamp *Bitstamp) buildPostForm(params *url.Values) {
	nonce := time.Now().UnixNano()
	//println(nonce)
	payload := fmt.Sprintf("%d%s%s", nonce, Bitstamp.clientId, Bitstamp.accessKey)
	sign, _ := GetParamHmacSHA256Sign(Bitstamp.secretkey, payload)
	params.Set("signature", strings.ToUpper(sign))
	params.Set("nonce", fmt.Sprintf("%d", nonce))
	params.Set("key", Bitstamp.accessKey)
}
func (Bitstamp *Bitstamp) GetAccount() (*Account, error) {
	urlStr := fmt.Sprintf("%s%s", BASE_URL, "v2/balance/")
	params := url.Values{}
	Bitstamp.buildPostForm(&params)
	resp, err := HttpPostForm(Bitstamp.client, urlStr, params)
	if err != nil {
		return nil, err
	}
	var respMap map[string]any
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return nil, err
	}
	acc := Account{}
	acc.Exchange = Bitstamp.String()
	acc.SubAccounts = make(map[Currency]SubAccount)
	acc.SubAccounts[BTC] = SubAccount{
		Currency:     BTC,
		Amount:       num.ToFloat64(respMap["btc_available"]),
		ForzenAmount: num.ToFloat64(respMap["btc_reserved"]),
		LoanAmount:   0,
	}
	acc.SubAccounts[LTC] = SubAccount{
		Currency:     LTC,
		Amount:       num.ToFloat64(respMap["ltc_available"]),
		ForzenAmount: num.ToFloat64(respMap["ltc_reserved"]),
		LoanAmount:   0,
	}
	acc.SubAccounts[ETH] = SubAccount{
		Currency:     ETH,
		Amount:       num.ToFloat64(respMap["eth_available"]),
		ForzenAmount: num.ToFloat64(respMap["eth_reserved"]),
		LoanAmount:   0,
	}
	acc.SubAccounts[XRP] = SubAccount{
		Currency:     XRP,
		Amount:       num.ToFloat64(respMap["xrp_available"]),
		ForzenAmount: num.ToFloat64(respMap["xrp_reserved"]),
		LoanAmount:   0,
	}
	acc.SubAccounts[USD] = SubAccount{
		Currency:     USD,
		Amount:       num.ToFloat64(respMap["usd_available"]),
		ForzenAmount: num.ToFloat64(respMap["usd_reserved"]),
		LoanAmount:   0,
	}
	acc.SubAccounts[EUR] = SubAccount{
		Currency:     EUR,
		Amount:       num.ToFloat64(respMap["eur_available"]),
		ForzenAmount: num.ToFloat64(respMap["eur_reserved"]),
		LoanAmount:   0,
	}
	acc.SubAccounts[BCH] = SubAccount{
		Currency:     BCH,
		Amount:       num.ToFloat64(respMap["bch_available"]),
		ForzenAmount: num.ToFloat64(respMap["bch_reserved"]),
		LoanAmount:   0}
	acc.SubAccounts[GBP] = SubAccount{
		Currency:     GBP,
		Amount:       num.ToFloat64(respMap["gbp_available"]),
		ForzenAmount: num.ToFloat64(respMap["gbp_reserved"]),
		LoanAmount:   0}
	acc.SubAccounts[PAX] = SubAccount{
		Currency:     PAX,
		Amount:       num.ToFloat64(respMap["pax_available"]),
		ForzenAmount: num.ToFloat64(respMap["pax_reserved"]),
		LoanAmount:   0}
	acc.SubAccounts[XLM] = SubAccount{
		Currency:     XLM,
		Amount:       num.ToFloat64(respMap["xlm_available"]),
		ForzenAmount: num.ToFloat64(respMap["xlm_reserved"]),
		LoanAmount:   0}
	return &acc, nil
}
func (Bitstamp *Bitstamp) placeOrder(side string, pair CurrencyPair, amount, price, urlStr string) (*Order, error) {
	params := url.Values{}
	params.Set("amount", amount)
	if price != "" {
		params.Set("price", price)
	}
	Bitstamp.buildPostForm(&params)
	resp, err := HttpPostForm(Bitstamp.client, urlStr, params)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return nil, err
	}
	orderId, isok := respMap["id"].(string)
	if !isok {
		return nil, errors.New(string(resp))
	}
	orderSide := BUY
	if side == "sell" {
		orderSide = SELL
	}
	orderprice, isok := respMap["price"].(string)
	if !isok {
		return nil, errors.New(string(resp))
	}
	return &Order{
		Currency:   pair,
		OrderID:    num.ToInt[int](orderId),
		OrderID2:   orderId,
		Price:      num.ToFloat64(orderprice),
		Amount:     num.ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       orderSide,
		Status:     ORDER_UNFINISH,
		OrderTime:  1}, nil
}
func (Bitstamp *Bitstamp) placeLimitOrder(side string, pair CurrencyPair, amount, price string) (*Order, error) {
	urlStr := fmt.Sprintf("%sv2/%s/%s/", BASE_URL, side, strings.ToLower(pair.ToSymbol("")))
	//println(urlStr)
	return Bitstamp.placeOrder(side, pair, amount, price, urlStr)
}
func (Bitstamp *Bitstamp) placeMarketOrder(side string, pair CurrencyPair, amount string) (*Order, error) {
	urlStr := fmt.Sprintf("%sv2/%s/market/%s/", BASE_URL, side, strings.ToLower(pair.ToSymbol("")))
	//println(urlStr)
	return Bitstamp.placeOrder(side, pair, amount, "", urlStr)
}
func (Bitstamp *Bitstamp) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return Bitstamp.placeLimitOrder("buy", currency, amount, price)
}
func (Bitstamp *Bitstamp) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return Bitstamp.placeLimitOrder("sell", currency, amount, price)
}
func (Bitstamp *Bitstamp) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return Bitstamp.placeMarketOrder("buy", currency, amount)
}
func (Bitstamp *Bitstamp) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return Bitstamp.placeMarketOrder("sell", currency, amount)
}
func (Bitstamp *Bitstamp) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	params := url.Values{}
	params.Set("id", orderId)
	Bitstamp.buildPostForm(&params)
	urlStr := BASE_URL + "v2/cancel_order/"
	resp, err := HttpPostForm(Bitstamp.client, urlStr, params)
	if err != nil {
		return false, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return false, err
	}
	if respMap["error"] != nil {
		return false, errors.New(string(resp))
	}
	println(string(resp))
	return true, nil
}
func (Bitstamp *Bitstamp) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	params := url.Values{}
	params.Set("id", orderId)
	Bitstamp.buildPostForm(&params)
	urlStr := BASE_URL + "order_status/"
	resp, err := HttpPostForm(Bitstamp.client, urlStr, params)
	if err != nil {
		return nil, err
	}
	//println(string(resp))
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return nil, err
	}
	transactions, isok := respMap["transactions"].([]any)
	if !isok {
		return nil, errors.New(string(resp))
	}
	status := respMap["status"].(string)
	ord := Order{}
	ord.Currency = currency
	ord.OrderID = num.ToInt[int](orderId)
	ord.OrderID2 = orderId
	if status == "Finished" {
		ord.Status = ORDER_FINISH
	} else {
		ord.Status = ORDER_UNFINISH
	}
	if len(transactions) > 0 {
		if ord.Status != ORDER_FINISH {
			ord.Status = ORDER_PART_FINISH
		}
		var (
			dealAmount  float64
			tradeAmount float64
			fee         float64
		)
		currencyStr := strings.ToLower(currency.CurrencyA.Symbol)
		for _, v := range transactions {
			transaction := v.(map[string]any)
			price := num.ToFloat64(transaction["price"])
			amount := num.ToFloat64(transaction[currencyStr])
			dealAmount += amount
			tradeAmount += amount * price
			fee += num.ToFloat64(transaction["fee"])
			//tpy := num.ToInt[int](transaction["type"]) //注意:不是交易方向，type (0 - deposit; 1 - withdrawal; 2 - market trade)
			//if tpy == 2 {
			//	ord.Side = SELL
			//}
		}
		avgPrice := tradeAmount / dealAmount
		ord.DealAmount = dealAmount
		ord.AvgPrice = avgPrice
		ord.Fee = fee
	}
	//	println(string(resp))
	return &ord, nil
}
func (Bitstamp *Bitstamp) GetUnfinishedOrders(currency CurrencyPair) ([]Order, error) {
	params := url.Values{}
	Bitstamp.buildPostForm(&params)
	urlStr := BASE_URL + "v2/open_orders/" + strings.ToLower(currency.ToSymbol("")) + "/"
	resp, err := HttpPostForm(Bitstamp.client, urlStr, params)
	if err != nil {
		return nil, err
	}
	respMap := make([]any, 1)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return nil, err
	}
	orders := make([]Order, 0)
	for _, v := range respMap {
		ord := v.(map[string]any)
		side := num.ToInt[int](ord["type"])
		orderSide := SELL
		if side == 0 {
			orderSide = BUY
		}
		orderTime, _ := time.Parse("2006-01-02 15:04:05", ord["datetime"].(string))
		orders = append(orders, Order{
			OrderID:   num.ToInt[int](ord["id"]),
			OrderID2:  fmt.Sprint(num.ToInt[int](ord["id"])),
			Currency:  currency,
			Price:     num.ToFloat64(ord["price"]),
			Amount:    num.ToFloat64(ord["amount"]),
			Side:      orderSide,
			Status:    ORDER_UNFINISH,
			OrderTime: int(orderTime.Unix())})
	}
	//println(string(resp))
	return orders, nil
}
func (Bitstamp *Bitstamp) GetOrderHistorys(currency CurrencyPair, optional ...OptionalParameter) ([]Order, error) {
	panic("not implement")
}

func (Bitstamp *Bitstamp) GetTicker(currency CurrencyPair) (*Ticker, error) {
	urlStr := BASE_URL + "v2/ticker/" + strings.ToLower(currency.ToSymbol(""))
	respMap, err := HttpGet(Bitstamp.client, urlStr)
	if err != nil {
		return nil, err
	}
	timestamp, _ := strconv.ParseUint(respMap["timestamp"].(string), 10, 64)
	return &Ticker{
		Pair: currency,
		Last: num.ToFloat64(respMap["last"]),
		High: num.ToFloat64(respMap["high"]),
		Low:  num.ToFloat64(respMap["low"]),
		Vol:  num.ToFloat64(respMap["volume"]),
		Sell: num.ToFloat64(respMap["ask"]),
		Buy:  num.ToFloat64(respMap["bid"]),
		Date: timestamp}, nil
}
func (Bitstamp *Bitstamp) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	urlStr := BASE_URL + "v2/order_book/" + strings.ToLower(currency.ToSymbol(""))
	respMap, err := HttpGet(Bitstamp.client, urlStr)
	if err != nil {
		return nil, err
	}
	//timestamp, _ := strconv.ParseUint(respMap["timestamp"].(string), 10, 64)
	bids, isok1 := respMap["bids"].([]any)
	asks, isok2 := respMap["asks"].([]any)
	if !isok1 || !isok2 {
		return nil, errors.New("get Depth Error")
	}
	i := 0
	dep := new(Depth)
	dep.Pair = currency
	for _, v := range bids {
		bid := v.([]any)
		dep.BidList = append(dep.BidList, DepthRecord{Price: num.ToFloat64(bid[0]), Amount: num.ToFloat64(bid[1])})
		i++
		if i == size {
			break
		}
	}
	i = 0
	for _, v := range asks {
		ask := v.([]any)
		dep.AskList = append(dep.AskList, DepthRecord{Price: num.ToFloat64(ask[0]), Amount: num.ToFloat64(ask[1])})
		i++
		if i == size {
			break
		}
	}
	sort.Sort(sort.Reverse(dep.AskList)) //reverse
	return dep, nil
}
func (Bitstamp *Bitstamp) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]Kline, error) {
	params := url.Values{}
	params.Set("step", _INTERNAL_KLINE_PERIOD_CONVERTER[period])
	params.Set("limit", fmt.Sprintf("%d", size))
	MergeOptionalParameter(&params, optional...)
	urlStr := BASE_URL + "v2/ohlc/" + strings.ToLower(currency.ToSymbol("")) + "?" + params.Encode()
	logs.D(urlStr)
	type ohlcResp struct {
		Data struct {
			Pair string `json:"pair"`
			Ohlc []struct {
				High      string `json:"high"`
				Timestamp string `json:"timestamp"`
				Volume    string `json:"volume"`
				Low       string `json:"low"`
				Close     string `json:"close"`
				Open      string `json:"open"`
			} `json:"ohlc"`
		} `json:"data"`
	}
	resp := ohlcResp{}
	err := HttpGet4(Bitstamp.client, urlStr, nil, &resp)
	if err != nil {
		return nil, err
	}
	var klineRecords []Kline
	for _, _record := range resp.Data.Ohlc {
		r := Kline{Pair: currency}
		r.Timestamp, _ = strconv.ParseInt(_record.Timestamp, 10, 64) //to unix timestramp
		r.Open = num.ToFloat64(_record.Open)
		r.High = num.ToFloat64(_record.High)
		r.Low = num.ToFloat64(_record.Low)
		r.Close = num.ToFloat64(_record.Close)
		r.Vol = num.ToFloat64(_record.Volume)
		klineRecords = append(klineRecords, r)
	}
	return klineRecords, nil
}

// //非个人，整个交易所的交易记录
func (Bitstamp *Bitstamp) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}
func (Bitstamp *Bitstamp) String() string {
	return BITSTAMP
}
