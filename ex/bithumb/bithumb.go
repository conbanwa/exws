package bithumb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/num"
	"log"
	"net/http"
	"net/url"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/web"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Bithumb struct {
	client *http.Client
	accesskey,
	secretkey string
}

var (
	baseUrl = "https://api.bithumb.com"
)

func New(client *http.Client, accesskey, secretkey string) *Bithumb {
	return &Bithumb{client: client, accesskey: accesskey, secretkey: secretkey}
}
func (bit *Bithumb) placeOrder(side, amount, price string, pair CurrencyPair) (*Order, error) {
	var retmap map[string]any
	params := fmt.Sprintf("order_currency=%s&units=%s&price=%s&type=%s", pair.CurrencyA.Symbol, amount, price, side)
	log.Println(params)
	err := bit.doAuthenticatedRequest("/trade/place", params, &retmap)
	if err != nil {
		return nil, err
	}
	if retmap["status"].(string) != "0000" {
		log.Println(retmap)
		return nil, errors.New(retmap["status"].(string))
	}
	var tradeSide TradeSide
	switch side {
	case "ask":
		tradeSide = SELL
	case "bid":
		tradeSide = BUY
	}
	log.Println(retmap)
	return &Order{
		OrderID:  num.ToInt[int](retmap["order_id"]),
		Amount:   num.ToFloat64(amount),
		Price:    num.ToFloat64(price),
		Currency: pair,
		Side:     tradeSide,
		Status:   ORDER_UNFINISH}, nil
}
func (bit *Bithumb) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return bit.placeOrder("bid", amount, price, currency)
}
func (bit *Bithumb) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return bit.placeOrder("ask", amount, price, currency)
}
func (bit *Bithumb) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bit *Bithumb) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implement")
}
func (bit *Bithumb) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("please invoke the CancelOrder2 method.")
}

/*补丁*/
func (bit *Bithumb) CancelOrder2(side, orderId string, currency CurrencyPair) (bool, error) {
	var retmap map[string]any
	params := fmt.Sprintf("type=%s&order_id=%s&currency=%s", side, orderId, currency.CurrencyA.Symbol)
	err := bit.doAuthenticatedRequest("/trade/cancel", params, &retmap)
	if err != nil {
		return false, err
	}
	if retmap["status"].(string) == "0000" {
		return true, nil
	}
	return false, errors.New(retmap["status"].(string))
}
func (bit *Bithumb) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("please invoke the GetOneOrder2 method.")
}

/*补丁*/
func (bit *Bithumb) GetOneOrder2(side, orderId string, currency CurrencyPair) (*Order, error) {
	var retmap map[string]any
	params := fmt.Sprintf("type=%s&order_id=%s&currency=%s", side, orderId, currency.CurrencyA.Symbol)
	err := bit.doAuthenticatedRequest("/info/order_detail", params, &retmap)
	if err != nil {
		return nil, err
	}
	if retmap["status"].(string) != "0000" {
		message := retmap["message"].(string)
		if "거래 체결내역이 존재하지 않습니다." == message {
			return nil, EX_ERR_NOT_FIND_ORDER
		}
		log.Println(retmap)
		return nil, errors.New(retmap["status"].(string))
	}
	order := new(Order)
	total := 0.0
	data := retmap["data"].([]any)
	for _, v := range data {
		ord := v.(map[string]any)
		switch ord["type"].(string) {
		case "ask":
			order.Side = SELL
		case "bid":
			order.Side = BUY
		}
		order.Amount += num.ToFloat64(ord["units_traded"])
		order.Fee += num.ToFloat64(ord["fee"])
		total += num.ToFloat64(ord["total"])
	}
	order.DealAmount = order.Amount
	avg := total / order.DealAmount
	order.AvgPrice = num.ToFloat64(fmt.Sprintf("%.2f", avg))
	order.Price = order.AvgPrice
	order.Currency = currency
	order.OrderID = num.ToInt[int](orderId)
	order.Status = ORDER_FINISH
	log.Println(retmap)
	return order, nil
}
func (bit *Bithumb) GetUnfinishedOrders(currency CurrencyPair) ([]Order, error) {
	var retmap map[string]any
	params := fmt.Sprintf("currency=%s", currency.CurrencyA.Symbol)
	err := bit.doAuthenticatedRequest("/info/orders", params, &retmap)
	if err != nil {
		return nil, err
	}
	if retmap["status"].(string) != "0000" {
		message := retmap["message"].(string)
		if "거래 진행중인 내역이 존재하지 않습니다." == message {
			return []Order{}, nil
		}
		return nil, fmt.Errorf("[%s]%s", retmap["status"].(string), message)
	}
	var orders []Order
	datas := retmap["data"].([]any)
	for _, v := range datas {
		orderinfo := v.(map[string]any)
		ord := Order{
			OrderID:  num.ToInt[int](orderinfo["order_id"]),
			Amount:   num.ToFloat64(orderinfo["units"]),
			Price:    num.ToFloat64(orderinfo["price"]),
			Currency: currency,
			Fee:      num.ToFloat64(orderinfo["fee"])}
		remaining := num.ToFloat64(orderinfo["units_remaining"])
		total := num.ToFloat64(orderinfo["total"])
		dealamount := ord.Amount - remaining
		ord.DealAmount = dealamount
		if dealamount > 0 {
			avg := fmt.Sprintf("%.4f", total/dealamount)
			ord.AvgPrice = num.ToFloat64(avg)
		}
		switch orderinfo["type"].(string) {
		case "ask":
			ord.Side = SELL
		case "bid":
			ord.Side = BUY
		}
		switch orderinfo["status"].(string) {
		case "placed":
			ord.Status = ORDER_UNFINISH
		}
		orders = append(orders, ord)
	}
	log.Println(retmap)
	return orders, nil
}
func (bit *Bithumb) GetOrderHistorys(currency CurrencyPair, optional ...OptionalParameter) ([]Order, error) {
	panic("not implement")
}
func (bit *Bithumb) GetAccount() (*Account, error) {
	var retmap map[string]any
	err := bit.doAuthenticatedRequest("/info/balance", "currency=ALL", &retmap)
	if err != nil {
		return nil, err
	}
	datamap := retmap["data"].(map[string]any)
	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount)
	for key := range datamap {
		if strings.HasPrefix(key, "available_") {
			t := strings.Split(key, "_")
			currency := NewCurrency(strings.ToUpper(t[len(t)-1]), "")
			acc.SubAccounts[currency] = SubAccount{
				Currency:     currency,
				Amount:       num.ToFloat64(datamap[key]),
				ForzenAmount: num.ToFloat64(datamap[fmt.Sprintf("in_use_%s", strings.ToLower(currency.String()))]),
				LoanAmount:   0}
		}
	}
	//log.Println(datamap)
	acc.Exchange = bit.String()
	return acc, nil
}
func (bit *Bithumb) doAuthenticatedRequest(uri, params string, ret any) error {
	nonce := time.Now().UnixNano() / int64(time.Millisecond)
	apiNonce := fmt.Sprint(nonce)
	eEndpoint := url.QueryEscape(uri)
	params += "&endpoint=" + eEndpoint
	// Api-Sign information generation.
	hmacData := uri + string(0) + params + string(0) + apiNonce
	hashHmacStr := GetParamHmacSHA512Base64Sign(bit.secretkey, hmacData)
	apiSign := hashHmacStr
	contentLengthStr := strconv.Itoa(len(params))
	// Connects to Bithumb API server and returns JSON result value.
	resp, err := web.NewRequest(bit.client, "POST", baseUrl+uri,
		bytes.NewBufferString(params).String(), map[string]string{
			"Api-Key":        bit.accesskey,
			"Api-Sign":       apiSign,
			"Api-Nonce":      apiNonce,
			"Content-Type":   "application/x-www-form-urlencoded",
			"Content-Length": contentLengthStr,
		}) // URL-encoded payload
	if err != nil {
		return err
	}
	err = json.Unmarshal(resp, ret)
	return err
}
func (bit *Bithumb) GetTicker(currency CurrencyPair) (*Ticker, error) {
	respMap, err := web.HttpGet(bit.client, fmt.Sprintf("%s/public/ticker/%s", baseUrl, currency.CurrencyA))
	if err != nil {
		return nil, err
	}
	s, isok := respMap["status"].(string)
	if s != "0000" || isok != true {
		msg := "ticker error"
		if isok {
			msg = s
		}
		return nil, errors.New(msg)
	}
	datamap := respMap["data"].(map[string]any)
	return &Ticker{
		Low:  num.ToFloat64(datamap["min_price"]),
		High: num.ToFloat64(datamap["max_price"]),
		Last: num.ToFloat64(datamap["closing_price"]),
		Vol:  num.ToFloat64(datamap["units_traded"]),
		Buy:  num.ToFloat64(datamap["buy_price"]),
		Sell: num.ToFloat64(datamap["sell_price"]),
	}, nil
}
func (bit *Bithumb) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	resp, err := web.HttpGet(bit.client, fmt.Sprintf("%s/public/orderbook/%s", baseUrl, currency.CurrencyA))
	if err != nil {
		return nil, err
	}
	if resp["status"].(string) != "0000" {
		return nil, errors.New(resp["status"].(string))
	}
	datamap := resp["data"].(map[string]any)
	bids := datamap["bids"].([]any)
	asks := datamap["asks"].([]any)
	dep := new(Depth)
	for _, v := range bids {
		bid := v.(map[string]any)
		dep.BidList = append(dep.BidList, DepthRecord{Price: num.ToFloat64(bid["price"]), Amount: num.ToFloat64(bid["quantity"])})
	}
	for _, v := range asks {
		ask := v.(map[string]any)
		dep.AskList = append(dep.AskList, DepthRecord{Price: num.ToFloat64(ask["price"]), Amount: num.ToFloat64(ask["quantity"])})
	}
	sort.Sort(sort.Reverse(dep.AskList))
	return dep, nil
}
func (bit *Bithumb) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, opt ...OptionalParameter) ([]Kline, error) {
	panic("not implement")
}

// 非个人，整个交易所的交易记录
func (bit *Bithumb) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implement")
}
func (bit *Bithumb) String() string {
	return BITHUMB
}
