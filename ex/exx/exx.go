package exx

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/num"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/web"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	EXX                       = "EXX"
	API_BASE_URL              = "https://api.exx.com/"
	MARKET_URL                = "http://api.exx.com/data/v1/"
	TICKER_API                = "ticker?currency=%s"
	DEPTH_API                 = "depth?currency=%s"
	TRADE_URL                 = "https://trade.exx.com/api/"
	GET_ACCOUNT_API           = "getBalance"
	GET_ORDER_API             = "getOrder"
	GET_UNFINISHED_ORDERS_API = "getUnfinishedOrdersIgnoreTradeType"
	CANCEL_ORDER_API          = "cancelOrder"
	PLACE_ORDER_API           = "order"
	WITHDRAW_API              = "withdraw"
	CANCELWITHDRAW_API        = "cancelWithdraw"
)

type Exx struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

func New(httpClient *http.Client, accessKey, secretKey string) *Exx {
	return &Exx{httpClient, accessKey, secretKey}
}
func (exx *Exx) String() string {
	return EXX
}
func (exx *Exx) GetTicker(currency cons.CurrencyPair) (*Ticker, error) {
	symbol := currency.ToLower().ToSymbol("_")
	path := MARKET_URL + fmt.Sprintf(TICKER_API, symbol)
	resp, err := web.HttpGet(exx.httpClient, path)
	if err != nil {
		return nil, err
	}
	//log.Println(path)
	//log.Println(resp)
	if err, isok := resp["error"].(string); isok {
		return nil, errors.New(err)
	}
	tickermap, _ := resp["ticker"].(map[string]any)
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
func (exx *Exx) GetDepth(size int, currency cons.CurrencyPair) (*Depth, error) {
	symbol := currency.ToSymbol("_")
	resp, err := web.HttpGet(exx.httpClient, MARKET_URL+fmt.Sprintf(DEPTH_API, symbol))
	if err != nil {
		return nil, err
	}
	log.Println(resp)
	//asks := resp["asks"].([]any)
	//bids := resp["bids"].([]any)
	asks, ok1 := resp["asks"].([]any)
	bids, ok2 := resp["bids"].([]any)
	if ok1 != true || ok2 != true {
		return nil, errors.New("no depth data")
	}
	log.Println(asks)
	log.Println(bids)
	depth := new(Depth)
	depth.Pair = currency
	for _, e := range bids {
		var r DepthRecord
		ee := e.([]any)
		r.Amount = num.ToFloat64(ee[1])
		r.Price = num.ToFloat64(ee[0])
		depth.BidList = append(depth.BidList, r)
	}
	for _, e := range asks {
		var r DepthRecord
		ee := e.([]any)
		r.Amount = num.ToFloat64(ee[1])
		r.Price = num.ToFloat64(ee[0])
		depth.AskList = append(depth.AskList, r)
	}
	sort.Sort(depth.AskList)
	depth.AskList = depth.AskList[0:size]
	depth.BidList = depth.BidList[0:size]
	return depth, nil
}
func (exx *Exx) buildPostForm(postForm *url.Values) error {
	postForm.Set("accesskey", exx.accessKey)
	nonce := time.Now().UnixNano() / 1000000
	postForm.Set("nonce", strconv.Itoa(int(nonce)))
	//postForm.Set("nonce", "1234567890123")
	payload := postForm.Encode()
	sign, err := GetParamHmacSHA512Sign(exx.secretKey, payload)
	if err != nil {
		return err
	}
	//log.Println(payload)
	postForm.Set("signature", sign)
	//postForm.Del("secret_key")
	//postForm.Set("reqTime", fmt.Sprintf("%d", time.Now().UnixNano()/1000000))
	return nil
}
func (exx *Exx) GetAccount() (*Account, error) {
	params := url.Values{}
	exx.buildPostForm(&params)
	log.Println(params.Encode())
	log.Println(TRADE_URL + GET_ACCOUNT_API + "?" + params.Encode())
	respMap, err := web.HttpGet(exx.httpClient, TRADE_URL+GET_ACCOUNT_API+"?"+params.Encode())
	if err != nil {
		return nil, err
	}
	if respMap["code"] != nil && respMap["code"].(float64) != 1000 {
		//return nil, errors.New(string(resp))
	}
	acc := new(Account)
	acc.Exchange = exx.String()
	acc.SubAccounts = make(map[cons.Currency]SubAccount)
	resultmap := respMap["result"].(map[string]any)
	coins := resultmap["coins"].([]any)
	acc.NetAsset = num.ToFloat64(resultmap["netAssets"])
	acc.Asset = num.ToFloat64(resultmap["totalAssets"])
	for _, v := range coins {
		vv := v.(map[string]any)
		subAcc := SubAccount{}
		subAcc.Amount = num.ToFloat64(vv["available"])
		subAcc.ForzenAmount = num.ToFloat64(vv["freez"])
		subAcc.Currency = cons.NewCurrency(vv["key"].(string), "").AdaptBchToBcc()
		acc.SubAccounts[subAcc.Currency] = subAcc
	}
	//log.Println(string(resp))
	//log.Println(acc)
	return acc, nil
}
func (exx *Exx) placeOrder(amount, price string, currency cons.CurrencyPair, tradeType int) (*q.Order, error) {
	symbol := currency.ToSymbol("_")
	params := url.Values{}
	params.Set("method", "order")
	params.Set("price", price)
	params.Set("amount", amount)
	params.Set("currency", symbol)
	params.Set("tradeType", fmt.Sprintf("%d", tradeType))
	exx.buildPostForm(&params)
	resp, err := web.HttpPostForm(exx.httpClient, TRADE_URL+PLACE_ORDER_API, params)
	if err != nil {
		//log.Println(err)
		return nil, err
	}
	//log.Println(string(resp));
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		//log.Println(err)
		return nil, err
	}
	code := respMap["code"].(float64)
	if code != 1000 {
		//log.Println(string(resp))
		return nil, fmt.Errorf("%.0f", code)
	}
	orid := respMap["id"].(string)
	order := new(q.Order)
	order.Amount, _ = strconv.ParseFloat(amount, 64)
	order.Price, _ = strconv.ParseFloat(price, 64)
	order.Status = cons.ORDER_UNFINISH
	order.Currency = currency
	order.OrderTime = int(time.Now().UnixNano() / 1000000)
	order.OrderID, _ = strconv.Atoi(orid)
	switch tradeType {
	case 0:
		order.Side = cons.SELL
	case 1:
		order.Side = cons.BUY
	}
	return order, nil
}
func (exx *Exx) LimitBuy(amount, price string, currency cons.CurrencyPair, opt ...cons.LimitOrderOptionalParameter) (*q.Order, error) {
	return exx.placeOrder(amount, price, currency, 1)
}
func (exx *Exx) LimitSell(amount, price string, currency cons.CurrencyPair, opt ...cons.LimitOrderOptionalParameter) (*q.Order, error) {
	return exx.placeOrder(amount, price, currency, 0)
}
func (exx *Exx) CancelOrder(orderId string, currency cons.CurrencyPair) (bool, error) {
	symbol := currency.ToSymbol("_")
	params := url.Values{}
	params.Set("method", "cancelOrder")
	params.Set("id", orderId)
	params.Set("currency", symbol)
	exx.buildPostForm(&params)
	resp, err := web.HttpPostForm(exx.httpClient, TRADE_URL+CANCEL_ORDER_API, params)
	if err != nil {
		//log.Println(err)
		return false, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		//log.Println(err)
		return false, err
	}
	code := respMap["code"].(float64)
	if code == 1000 {
		return true, nil
	}
	//log.Println(respMap)
	return false, fmt.Errorf("%.0f", code)
}
func parseOrder(order *q.Order, ordermap map[string]any) {
	//log.Println(ordermap)
	//order.Currency = currency;
	order.OrderID, _ = strconv.Atoi(ordermap["id"].(string))
	order.OrderID2 = ordermap["id"].(string)
	order.Amount = ordermap["total_amount"].(float64)
	order.DealAmount = ordermap["trade_amount"].(float64)
	order.Price = ordermap["price"].(float64)
	//	order.Fee = ordermap["fees"].(float64)
	if order.DealAmount > 0 {
		order.AvgPrice = ordermap["trade_money"].(float64) / order.DealAmount
	} else {
		order.AvgPrice = 0
	}
	order.OrderTime = int(ordermap["trade_date"].(float64))
	orType := ordermap["type"].(float64)
	switch orType {
	case 0:
		order.Side = cons.SELL
	case 1:
		order.Side = cons.BUY
	default:
		log.Printf("unknown order type %f", orType)
	}
	_status := cons.TradeStatus(ordermap["status"].(float64))
	switch _status {
	case 0:
		order.Status = cons.ORDER_UNFINISH
	case 1:
		order.Status = cons.ORDER_CANCEL
	case 2:
		order.Status = cons.ORDER_FINISH
	case 3:
		order.Status = cons.ORDER_UNFINISH
	}
}
func (exx *Exx) GetOneOrder(orderId string, currency cons.CurrencyPair) (*q.Order, error) {
	symbol := currency.ToSymbol("_")
	params := url.Values{}
	params.Set("method", "getOrder")
	params.Set("id", orderId)
	params.Set("currency", symbol)
	exx.buildPostForm(&params)
	resp, err := web.HttpPostForm(exx.httpClient, TRADE_URL+GET_ORDER_API, params)
	if err != nil {
		//log.Println(err)
		return nil, err
	}
	//println(string(resp))
	ordermap := make(map[string]any)
	err = json.Unmarshal(resp, &ordermap)
	if err != nil {
		//log.Println(err)
		return nil, err
	}
	order := new(q.Order)
	order.Currency = currency
	parseOrder(order, ordermap)
	return order, nil
}
func (exx *Exx) GetUnfinishedOrders(currency cons.CurrencyPair) ([]q.Order, error) {
	params := url.Values{}
	symbol := currency.ToSymbol("_")
	params.Set("method", "getUnfinishedOrdersIgnoreTradeType")
	params.Set("currency", symbol)
	params.Set("pageIndex", "1")
	params.Set("pageSize", "100")
	exx.buildPostForm(&params)
	resp, err := web.HttpPostForm(exx.httpClient, TRADE_URL+GET_UNFINISHED_ORDERS_API, params)
	if err != nil {
		//log.Println(err)
		return nil, err
	}
	respstr := string(resp)
	//println(respstr)
	if strings.Contains(respstr, "\"code\":3001") {
		//log.Println(respstr)
		return nil, nil
	}
	var resps []any
	err = json.Unmarshal(resp, &resps)
	if err != nil {
		//log.Println(err)
		return nil, err
	}
	var orders []q.Order
	for _, v := range resps {
		ordermap := v.(map[string]any)
		order := q.Order{}
		order.Currency = currency
		parseOrder(&order, ordermap)
		orders = append(orders, order)
	}
	return orders, nil
}
func (exx *Exx) GetOrderHistorys(currency cons.CurrencyPair, optional ...OptionalParameter) ([]q.Order, error) {
	return nil, nil
}
func (exx *Exx) GetKlineRecords(currency cons.CurrencyPair, period, size, since int) ([]Kline, error) {
	return nil, nil
}
func (exx *Exx) Withdraw(amount string, currency cons.Currency, fees, receiveAddr, safePwd string) (string, error) {
	params := url.Values{}
	params.Set("method", "withdraw")
	params.Set("currency", strings.ToLower(currency.AdaptBchToBcc().String()))
	params.Set("amount", amount)
	params.Set("fees", fees)
	params.Set("receiveAddr", receiveAddr)
	params.Set("safePwd", safePwd)
	exx.buildPostForm(&params)
	resp, err := web.HttpPostForm(exx.httpClient, TRADE_URL+WITHDRAW_API, params)
	if err != nil {
		//log.Println("withdraw fail.", err)
		return "", err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		//log.Println(err, string(resp))
		return "", err
	}
	if respMap["code"].(float64) == 1000 {
		return respMap["id"].(string), nil
	}
	return "", errors.New(string(resp))
}
func (exx *Exx) CancelWithdraw(id string, currency cons.Currency, safePwd string) (bool, error) {
	params := url.Values{}
	params.Set("method", "cancelWithdraw")
	params.Set("currency", strings.ToLower(currency.AdaptBchToBcc().String()))
	params.Set("downloadId", id)
	params.Set("safePwd", safePwd)
	exx.buildPostForm(&params)
	resp, err := web.HttpPostForm(exx.httpClient, TRADE_URL+CANCELWITHDRAW_API, params)
	if err != nil {
		//log.Println("cancel withdraw fail.", err)
		return false, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		//log.Println(err, string(resp))
		return false, err
	}
	if respMap["code"].(float64) == 1000 {
		return true, nil
	}
	return false, errors.New(string(resp))
}
func (exx *Exx) GetTrades(currencyPair cons.CurrencyPair, since int64) ([]q.Trade, error) {
	panic("unimplements")
}
func (exx *Exx) MarketBuy(amount, price string, currency cons.CurrencyPair) (*q.Order, error) {
	panic("unsupport the market order")
}
func (exx *Exx) MarketSell(amount, price string, currency cons.CurrencyPair) (*q.Order, error) {
	panic("unsupport the market order")
}
