package allcoin

import (
	"encoding/json"
	"errors"
	"github.com/conbanwa/num"
	"github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/web"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	API_BASE_URL           = "https://www.allcoin.ca/"
	TICKER_URI             = "Api_Market/getCoinTrade"
	TICKER_URI_2           = "Api_Order/ticker"
	TICKERS_URI            = "ticker/allBookTickers"
	DEPTH_URI              = "Api_Order/depth"
	ACCOUNT_URI            = "Api_User/userBalance"
	ORDER_URI              = "Api_Order/coinTrust"
	ORDER_CANCEL_URI       = "Api_Order/cancel"
	ORDER_INFO_URI         = "Api_Order/orderInfo"
	UNFINISHED_ORDERS_INFO = "Api_Order/trustList"
)

type Allcoin struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func (ac *Allcoin) buildParamsSigned(postForm *url.Values) error {
	//postForm.Set("api_key", ac.accessKey)
	//postForm.Set("secret_key", ac.secretKey)
	payload := postForm.Encode() + "&secret_key=" + ac.secretKey
	//log.Println("payload:", payload, "postForm:", postForm.Encode())
	sign, _ := wstrader.GetParamMD5Sign(ac.secretKey, payload)
	postForm.Set("sign", sign)
	return nil
}
func New(client *http.Client, apiKey, secretKey string) *Allcoin {
	return &Allcoin{apiKey, secretKey, client}
}
func (ac *Allcoin) String() string {
	return "allcoin.com"
}
func (ac *Allcoin) GetTicker(currency cons.CurrencyPair) (*wstrader.Ticker, error) {
	//wg := sync.WaitGroup{}
	//wg.Add(2)
	//go func() {
	//	defer wg.Done()
	//currency2 := ac.adaptCurrencyPair(currency)
	//params := url.Values{}
	//params.Set("symbol", strings.ToLower(currency2.ToSymbol("2")))
	//path := API_BASE_URL + TICKER_URI_2
	//resp, err := HttpPostForm(ac.httpClient, path, params)
	//log.Println("resp:", string(resp), "err:", err)
	//if err != nil {
	//	//return nil, err
	//}
	//
	//respMap := make(map[string]any)
	//err = json.Unmarshal(resp, &respMap)
	//if err != nil {
	//	log.Println(string(resp))
	//	//return nil, err
	//}
	//code := respMap["code"].(float64)
	//msg := respMap["msg"].(string)
	//log.Println("code=", code, "msg:", msg)
	//if code != 0 {
	//	//return nil, errors.New(respMap["msg"].(string))
	//}
	//data := respMap["data"].(map[string]any)
	//log.Println("1", data)
	//}()
	//go func() {
	//	defer wg.Done()
	currency2 := ac.adaptCurrencyPair(currency)
	params := url.Values{}
	params.Set("part", strings.ToLower(currency2.CurrencyB.String()))
	params.Set("coin", strings.ToLower(currency2.CurrencyA.String()))
	path := API_BASE_URL + TICKER_URI
	resp, err := web.HttpPostForm(ac.httpClient, path, params)
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}
	//log.Println("2", respMap)
	//}()
	//wg.Wait()
	var ticker wstrader.Ticker
	ticker.Pair = currency
	ticker.Date = uint64(time.Now().Unix())
	ticker.Last = num.ToFloat64(respMap["price"])
	ticker.Buy = num.ToFloat64(respMap["buy"])
	ticker.Sell = num.ToFloat64(respMap["sale"])
	ticker.Low = num.ToFloat64(respMap["min"])
	ticker.High = num.ToFloat64(respMap["max"])
	ticker.Vol = num.ToFloat64(respMap["volume_24h"])
	return &ticker, nil
}
func (ac *Allcoin) GetDepth(size int, currencyPair cons.CurrencyPair) (*wstrader.Depth, error) {
	currency2 := ac.adaptCurrencyPair(currencyPair)
	params := url.Values{}
	params.Set("symbol", strings.ToLower(currency2.ToSymbol("2")))
	path := API_BASE_URL + DEPTH_URI
	resp, err := web.HttpPostForm(ac.httpClient, path, params)
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}
	code := respMap["code"].(float64)
	msg := respMap["msg"].(string)
	log.Println("code=", code, "msg:", msg)
	if code != 0 {
		return nil, errors.New(respMap["msg"].(string))
	}
	data := respMap["data"].(map[string]any)
	log.Println("1", data)
	bids := data["bids"].([]any)
	asks := data["asks"].([]any)
	//log.Println("len bids", len(bids))
	//log.Println("len asks", len(asks))
	depth := new(wstrader.Depth)
	depth.Pair = currencyPair
	for _, bid := range bids {
		_bid := bid.([]any)
		amount := num.ToFloat64(_bid[1])
		price := num.ToFloat64(_bid[0])
		dr := wstrader.DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}
	for _, ask := range asks {
		_ask := ask.([]any)
		amount := num.ToFloat64(_ask[1])
		price := num.ToFloat64(_ask[0])
		dr := wstrader.DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}
	sort.Sort(depth.AskList)
	return depth, nil
}
func (ac *Allcoin) placeOrder(amount, price string, pair cons.CurrencyPair, orderType, orderSide string) (*q.Order, error) {
	pair = ac.adaptCurrencyPair(pair)
	path := API_BASE_URL + ORDER_URI
	params := url.Values{}
	params.Set("api_key", ac.accessKey)
	params.Set("symbol", strings.ToLower(pair.ToSymbol("2")))
	params.Set("type", orderSide)
	params.Set("price", price)
	params.Set("number", amount)
	ac.buildParamsSigned(&params)
	resp, err := web.HttpPostForm(ac.httpClient, path, params)
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}
	code := respMap["code"].(float64)
	if code != 0 {
		return nil, errors.New(respMap["msg"].(string))
	}
	data := respMap["data"].(map[string]any)
	orderId := data["order_id"].(string)
	side := cons.BUY
	if orderSide == "sale" {
		side = cons.SELL
	}
	return &q.Order{
		Currency: pair,
		//OrderID:
		OrderID2:   orderId,
		Price:      num.ToFloat64(price),
		Amount:     num.ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       side,
		Status:     cons.ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}
func (ac *Allcoin) GetAccount() (*wstrader.Account, error) {
	params := url.Values{}
	params.Set("api_key", ac.accessKey)
	ac.buildParamsSigned(&params)
	//log.Println("params=", params)
	path := API_BASE_URL + ACCOUNT_URI
	resp, err := web.HttpPostForm(ac.httpClient, path, params)
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}
	code := respMap["code"].(float64)
	//msg := respMap["msg"].(string)
	//log.Println("code=", code, "msg:", msg)
	if code != 0 {
		return nil, errors.New(respMap["msg"].(string))
	}
	data := respMap["data"].(map[string]any)
	acc := wstrader.Account{}
	acc.Exchange = ac.String()
	acc.SubAccounts = make(map[cons.Currency]wstrader.SubAccount)
	for k, v := range data {
		s := strings.Split(k, "_")
		if len(s) == 2 {
			cur := cons.NewCurrency(s[0], "")
			if s[1] == "over" {
				sub := wstrader.SubAccount{}
				sub = acc.SubAccounts[cur]
				sub.Amount = num.ToFloat64(v)
				acc.SubAccounts[cur] = sub
			} else if s[1] == "lock" {
				sub := wstrader.SubAccount{}
				sub = acc.SubAccounts[cur]
				sub.ForzenAmount = num.ToFloat64(v)
				acc.SubAccounts[cur] = sub
			}
		}
	}
	log.Println(acc)
	return &acc, nil
}
func (ac *Allcoin) LimitBuy(amount, price string, currencyPair cons.CurrencyPair, opt ...cons.LimitOrderOptionalParameter) (*q.Order, error) {
	return ac.placeOrder(amount, price, currencyPair, "LIMIT", "buy")
}
func (ac *Allcoin) LimitSell(amount, price string, currencyPair cons.CurrencyPair, opt ...cons.LimitOrderOptionalParameter) (*q.Order, error) {
	return ac.placeOrder(amount, price, currencyPair, "LIMIT", "sale")
}
func (ac *Allcoin) MarketBuy(amount, price string, currencyPair cons.CurrencyPair) (*q.Order, error) {
	return ac.placeOrder(amount, price, currencyPair, "MARKET", "buy")
}
func (ac *Allcoin) MarketSell(amount, price string, currencyPair cons.CurrencyPair) (*q.Order, error) {
	return ac.placeOrder(amount, price, currencyPair, "MARKET", "sale")
}
func (ac *Allcoin) CancelOrder(orderId string, currencyPair cons.CurrencyPair) (bool, error) {
	currencyPair = ac.adaptCurrencyPair(currencyPair)
	path := API_BASE_URL + ORDER_CANCEL_URI
	params := url.Values{}
	params.Set("api_key", ac.accessKey)
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("2")))
	params.Set("order_id", orderId)
	ac.buildParamsSigned(&params)
	resp, err := web.HttpPostForm(ac.httpClient, path, params)
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return false, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return false, err
	}
	code := respMap["code"].(int)
	if code != 0 {
		return false, errors.New(respMap["msg"].(string))
	}
	//orderIdCanceled := num.ToInt[int](respMap["orderId"])
	//if orderIdCanceled <= 0 {
	//	return false, errors.New(string(resp))
	//}
	return true, nil
}
func (ac *Allcoin) GetOneOrder(orderId string, currencyPair cons.CurrencyPair) (*q.Order, error) {
	currencyPair = ac.adaptCurrencyPair(currencyPair)
	path := API_BASE_URL + ORDER_INFO_URI
	params := url.Values{}
	params.Set("api_key", ac.accessKey)
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("2")))
	params.Set("trust_id", orderId)
	err := ac.buildParamsSigned(&params)
	if err != nil {
		return nil, err
	}
	resp, err := web.HttpPostForm(ac.httpClient, path, params)
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}
	code := respMap["code"].(float64)
	if code != 0 {
		return nil, errors.New(respMap["msg"].(string))
	}
	data := respMap["data"].(map[string]any)
	status := data["status"]
	side := data["flag"]
	ord := q.Order{}
	ord.Currency = currencyPair
	//ord.OrderID = num.ToInt[int](orderId)
	ord.OrderID2 = orderId
	if side == "sale" {
		ord.Side = cons.SELL
	} else {
		ord.Side = cons.BUY
	}
	switch status {
	case "1": //TODO
		ord.Status = cons.ORDER_FINISH
	case "2":
		ord.Status = cons.ORDER_PART_FINISH
	case "3":
		ord.Status = cons.ORDER_CANCEL
	case "PENDING_CANCEL":
		ord.Status = cons.ORDER_CANCEL_ING
	case "REJECTED":
		ord.Status = cons.ORDER_REJECT
	}
	ord.Amount = num.ToFloat64(data["number"])
	ord.Price = num.ToFloat64(data["price"])
	ord.DealAmount = ord.Amount - num.ToFloat64(data["numberover"])
	ord.AvgPrice = num.ToFloat64(data["avg_price"]) // response no avg price ， fill price
	return &ord, nil
}
func (ac *Allcoin) GetUnfinishedOrders(currencyPair cons.CurrencyPair) ([]q.Order, error) {
	currencyPair = ac.adaptCurrencyPair(currencyPair)
	path := API_BASE_URL + UNFINISHED_ORDERS_INFO
	params := url.Values{}
	params.Set("api_key", ac.accessKey)
	params.Set("symbol", strings.ToLower(currencyPair.ToSymbol("2")))
	params.Set("type", "open")
	err := ac.buildParamsSigned(&params)
	if err != nil {
		return nil, err
	}
	resp, err := web.HttpPostForm(ac.httpClient, path, params)
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}
	code := respMap["code"].(float64)
	//msg := respMap["msg"].(string)
	//log.Println("code=", code, "msg:", msg)
	if code != 0 {
		return nil, errors.New(respMap["msg"].(string))
	}
	data, isok := respMap["data"].([]map[string]any)
	orders := make([]q.Order, 0)
	if isok != true {
		return orders, nil
	}
	for _, ord := range data {
		//ord := v.(map[string]any)
		//side := ord["side"].(string)
		//orderSide := SELL
		//if side == "BUY" {
		//	orderSide = BUY
		//}
		orders = append(orders, q.Order{
			OrderID:  num.ToInt[int](ord["id"]),
			OrderID2: ord["id"].(string),
			Currency: currencyPair,
			Price:    num.ToFloat64(ord["price"]),
			Amount:   num.ToFloat64(ord["number"]),
			//Side:      TradeSide(orderSide),
			//Status:    ORDER_UNFINISH,
			OrderTime: num.ToInt[int](ord["created"])})
	}
	return orders, nil
}
func (ac *Allcoin) GetKlineRecords(currency cons.CurrencyPair, period, size, since int) ([]wstrader.Kline, error) {
	panic("not implements")
}

// 非个人，整个交易所的交易记录
func (ac *Allcoin) GetTrades(currencyPair cons.CurrencyPair, since int64) ([]q.Trade, error) {
	panic("not implements")
}
func (ac *Allcoin) GetOrderHistorys(currency cons.CurrencyPair, opt ...wstrader.OptionalParameter) ([]q.Order, error) {
	panic("not implements")
}
func (ac *Allcoin) adaptCurrencyPair(pair cons.CurrencyPair) cons.CurrencyPair {
	return pair
}
