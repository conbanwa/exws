package zb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/logs"
	"github.com/conbanwa/num"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	MarketUrl              = "http://api.zb.com/data/v1/"
	TickerApi              = "ticker?market=%s"
	DepthApi               = "depth?market=%s&size=%d"
	TradeUrl               = "https://trade.zb.com/api/"
	GetAccountApi          = "getAccountInfo"
	GetOrderApi            = "getOrder"
	GetUnfinishedOrdersApi = "getUnfinishedOrdersIgnoreTradeType"
	CancelOrderApi         = "cancelOrder"
	PlaceOrderApi          = "order"
	WithdrawApi            = "withdraw"
	CancelwithdrawApi      = "cancelWithdraw"
)

type Zb struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

func New(httpClient *http.Client, accessKey, secretKey string) *Zb {
	return &Zb{httpClient, accessKey, secretKey}
}
func (zb *Zb) String() string {
	return ZB
}
func (zb *Zb) GetTicker(currency CurrencyPair) (*Ticker, error) {
	symbol := currency.ToSymbol("_")
	resp, err := HttpGet(zb.httpClient, MarketUrl+fmt.Sprintf(TickerApi, symbol))
	if err != nil {
		return nil, err
	}
	//logs.E(resp)
	tickermap := resp["ticker"].(map[string]any)
	ticker := new(Ticker)
	ticker.Pair = currency
	ticker.Date, _ = strconv.ParseUint(resp["date"].(string), 10, 64)
	ticker.Buy, _ = strconv.ParseFloat(tickermap["buy"].(string), 64)
	ticker.Sell, _ = strconv.ParseFloat(tickermap["sell"].(string), 64)
	ticker.Last, _ = strconv.ParseFloat(tickermap["last"].(string), 64)
	ticker.High, _ = strconv.ParseFloat(tickermap["high"].(string), 64)
	ticker.Low, _ = strconv.ParseFloat(tickermap["low"].(string), 64)
	ticker.Vol, _ = strconv.ParseFloat(tickermap["vol"].(string), 64)
	return ticker, nil
}
func (zb *Zb) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	symbol := currency.ToSymbol("_")
	resp, err := HttpGet(zb.httpClient, MarketUrl+fmt.Sprintf(DepthApi, symbol, size))
	if err != nil {
		return nil, err
	}
	asks, isok1 := resp["asks"].([]any)
	bids, isok2 := resp["bids"].([]any)
	if !isok2 || !isok1 {
		return nil, errors.New("no depth data")
	}
	//logs.E(asks)
	//logs.E(bids)
	depth := new(Depth)
	depth.Pair = currency
	for _, e := range bids {
		var r DepthRecord
		ee := e.([]any)
		r.Amount = ee[1].(float64)
		r.Price = ee[0].(float64)
		depth.BidList = append(depth.BidList, r)
	}
	for _, e := range asks {
		var r DepthRecord
		ee := e.([]any)
		r.Amount = ee[1].(float64)
		r.Price = ee[0].(float64)
		depth.AskList = append(depth.AskList, r)
	}
	return depth, nil
}
func (zb *Zb) buildPostForm(postForm *url.Values) error {
	postForm.Set("accesskey", zb.accessKey)
	payload := postForm.Encode()
	secretkeySha, _ := GetSHA(zb.secretKey)
	sign, err := GetParamHmacMD5Sign(secretkeySha, payload)
	if err != nil {
		return err
	}
	postForm.Set("sign", sign)
	//postForm.Del("secret_key")
	postForm.Set("reqTime", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	return nil
}
func (zb *Zb) GetAccount() (*Account, error) {
	params := url.Values{}
	params.Set("method", "getAccountInfo")
	zb.buildPostForm(&params)
	//logs.E(params.Encode())
	resp, err := HttpPostForm(zb.httpClient, TradeUrl+GetAccountApi, params)
	if err != nil {
		return nil, err
	}
	var respMap map[string]any
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		logs.E("json unmarshal error")
		return nil, err
	}
	if respMap["code"] != nil && respMap["code"].(float64) != 1000 {
		return nil, errors.New(string(resp))
	}
	acc := new(Account)
	acc.Exchange = zb.String()
	acc.SubAccounts = make(map[Currency]SubAccount)
	resultmap := respMap["result"].(map[string]any)
	coins := resultmap["coins"].([]any)
	acc.NetAsset = num.ToFloat64(resultmap["netAssets"])
	acc.Asset = num.ToFloat64(resultmap["totalAssets"])
	for _, v := range coins {
		vv := v.(map[string]any)
		subAcc := SubAccount{}
		subAcc.Amount = num.ToFloat64(vv["available"])
		subAcc.ForzenAmount = num.ToFloat64(vv["freez"])
		subAcc.Currency = NewCurrency(vv["key"].(string), "").AdaptBchToBcc()
		acc.SubAccounts[subAcc.Currency] = subAcc
	}
	//logs.E(string(resp))
	//logs.E(acc)
	return acc, nil
}
func (zb *Zb) placeOrder(amount, price string, currency CurrencyPair, tradeType int) (*Order, error) {
	symbol := currency.ToSymbol("_")
	params := url.Values{}
	params.Set("method", "order")
	params.Set("price", price)
	params.Set("amount", amount)
	params.Set("currency", symbol)
	params.Set("tradeType", fmt.Sprintf("%d", tradeType))
	zb.buildPostForm(&params)
	resp, err := HttpPostForm(zb.httpClient, TradeUrl+PlaceOrderApi, params)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	//logs.E(string(resp));
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	if code := respMap["code"].(float64); code != 1000 {
		logs.E(string(resp))
		return nil, fmt.Errorf("%.0f", code)
	}
	orid := respMap["id"].(string)
	order := new(Order)
	order.Amount, _ = strconv.ParseFloat(amount, 64)
	order.Price, _ = strconv.ParseFloat(price, 64)
	order.Status = ORDER_UNFINISH
	order.Currency = currency
	order.OrderTime = int(time.Now().UnixNano() / 1000000)
	order.OrderID, _ = strconv.Atoi(orid)
	switch tradeType {
	case 0:
		order.Side = SELL
	case 1:
		order.Side = BUY
	}
	return order, nil
}
func (zb *Zb) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return zb.placeOrder(amount, price, currency, 1)
}
func (zb *Zb) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return zb.placeOrder(amount, price, currency, 0)
}
func (zb *Zb) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	symbol := currency.ToSymbol("-")
	params := url.Values{}
	params.Set("method", "cancelOrder")
	params.Set("id", orderId)
	params.Set("currency", symbol)
	zb.buildPostForm(&params)
	resp, err := HttpPostForm(zb.httpClient, TradeUrl+CancelOrderApi, params)
	if err != nil {
		logs.E(err)
		return false, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		logs.E(err)
		return false, err
	}
	code := respMap["code"].(float64)
	if code == 1000 {
		return true, nil
	}
	//logs.E(respMap)
	return false, fmt.Errorf("%.0f", code)
}
func parseOrder(order *Order, ordermap map[string]any) {
	//logs.E(ordermap)
	//order.Currency = currency;
	order.OrderID, _ = strconv.Atoi(ordermap["id"].(string))
	order.OrderID2 = ordermap["id"].(string)
	order.Amount = ordermap["total_amount"].(float64)
	order.DealAmount = ordermap["trade_amount"].(float64)
	order.Price = ordermap["price"].(float64)
	//	order.Fee = ordermap["fees"].(float64)
	if order.DealAmount > 0 {
		order.AvgPrice = num.ToFloat64(ordermap["trade_money"]) / order.DealAmount
	} else {
		order.AvgPrice = 0
	}
	order.OrderTime = int(ordermap["trade_date"].(float64))
	orType := ordermap["type"].(float64)
	switch orType {
	case 0:
		order.Side = SELL
	case 1:
		order.Side = BUY
	default:
		logs.E("unknown order type %f", orType)
	}
	_status := TradeStatus(ordermap["status"].(float64))
	switch _status {
	case 0:
		order.Status = ORDER_UNFINISH
	case 1:
		order.Status = ORDER_CANCEL
	case 2:
		order.Status = ORDER_FINISH
	case 3:
		order.Status = ORDER_UNFINISH
	}
}
func (zb *Zb) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	symbol := currency.ToSymbol("_")
	params := url.Values{}
	params.Set("method", "getOrder")
	params.Set("id", orderId)
	params.Set("currency", symbol)
	zb.buildPostForm(&params)
	resp, err := HttpPostForm(zb.httpClient, TradeUrl+GetOrderApi, params)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	//println(string(resp))
	ordermap := make(map[string]any)
	err = json.Unmarshal(resp, &ordermap)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	order := new(Order)
	order.Currency = currency
	parseOrder(order, ordermap)
	return order, nil
}
func (zb *Zb) GetUnfinishedOrders(currency CurrencyPair) ([]Order, error) {
	params := url.Values{}
	symbol := currency.ToSymbol("_")
	params.Set("method", "getUnfinishedOrdersIgnoreTradeType")
	params.Set("currency", symbol)
	params.Set("pageIndex", "1")
	params.Set("pageSize", "100")
	zb.buildPostForm(&params)
	resp, err := HttpPostForm(zb.httpClient, TradeUrl+GetUnfinishedOrdersApi, params)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	respstr := string(resp)
	//println(respstr)
	if strings.Contains(respstr, "\"code\":3001") {
		logs.E(respstr)
		return nil, nil
	}
	var resps []any
	err = json.Unmarshal(resp, &resps)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	var orders []Order
	for _, v := range resps {
		ordermap := v.(map[string]any)
		order := Order{}
		order.Currency = currency
		parseOrder(&order, ordermap)
		orders = append(orders, order)
	}
	return orders, nil
}
func (zb *Zb) GetOrderHistorys(currency CurrencyPair, opt ...OptionalParameter) ([]Order, error) {
	return nil, nil
}
func (zb *Zb) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, opt ...OptionalParameter) ([]Kline, error) {
	return nil, nil
}
func (zb *Zb) Withdraw(amount string, currency Currency, fees, receiveAddr, safePwd string) (string, error) {
	params := url.Values{}
	params.Set("method", "withdraw")
	params.Set("currency", strings.ToLower(currency.AdaptBchToBcc().String()))
	params.Set("amount", amount)
	params.Set("fees", fees)
	params.Set("receiveAddr", receiveAddr)
	params.Set("safePwd", safePwd)
	zb.buildPostForm(&params)
	resp, err := HttpPostForm(zb.httpClient, TradeUrl+WithdrawApi, params)
	if err != nil {
		logs.E("withdraw fail.", err)
		return "", err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		logs.E(err, string(resp))
		return "", err
	}
	if respMap["code"].(float64) == 1000 {
		return respMap["id"].(string), nil
	}
	return "", errors.New(string(resp))
}
func (zb *Zb) CancelWithdraw(id string, currency Currency, safePwd string) (bool, error) {
	params := url.Values{}
	params.Set("method", "cancelWithdraw")
	params.Set("currency", strings.ToLower(currency.AdaptBchToBcc().String()))
	params.Set("downloadId", id)
	params.Set("safePwd", safePwd)
	zb.buildPostForm(&params)
	resp, err := HttpPostForm(zb.httpClient, TradeUrl+CancelwithdrawApi, params)
	if err != nil {
		logs.E("cancel withdraw fail.", err)
		return false, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		logs.E(err, string(resp))
		return false, err
	}
	if respMap["code"].(float64) == 1000 {
		return true, nil
	}
	return false, errors.New(string(resp))
}
func (zb *Zb) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("unimplements")
}
func (zb *Zb) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unsupport the market order")
}
func (zb *Zb) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unsupport the market order")
}
