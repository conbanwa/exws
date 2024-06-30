package atop

/**
CancelOrder(orderId string, currency CurrencyPair) (bool, error)
GetOneOrder(orderId string, currency CurrencyPair) (*Order, error)
GetUnfinishedOrders(currency CurrencyPair) ([]Order, error)
GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error)
GetAccount() (*Account, error)
GetDepth(size int, currency CurrencyPair) (*Depth, error)
GetKlineRecords(currency CurrencyPair, period , size, since int) ([]Kline, error)
//Non-individual, transaction record of the entire exchange
GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error)
String() string
*/
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/num"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	. "github.com/conbanwa/wstrader/util"
	. "github.com/conbanwa/wstrader/web"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	//test  https://testapi.a.top
	//product   https://api.a.top
	ApiBaseUrl = "https://testapi.a.top"
	//ApiBaseUrl = "https://api.a.top"
	////market data
	//Trading market configuration
	GetMarketConfig = "/data/api/v1/getMarketConfig"
	//K line data
	GetKLine = "/data/api/v1/getKLine"
	//Aggregate market
	GetTicker = "/data/api/v1/getTicker?market=%s"
	//The latest Ticker for all markets
	GetTickers = "/data/api/v1/getTickers"
	//Market depth data
	GetDepth = "/data/api/v1/getDepth?market=%s"
	//Recent market record
	GetTrades = "/data/api/v1/getTrades"
	////trading
	//Get server time (no signature required)
	GetServerTime = "/trade/api/v1/getServerTime"
	//Get atcount balance
	GetBalance = "/trade/api/v1/getBalance"
	//Plate the order
	PlateOrder = "/trade/api/v1/order"
	//Commissioned by batch
	BatchOrder = "/trade/api/v1/batchOrder"
	//cancellations
	CancelOrder = "/trade/api/v1/cancel"
	//From a single batch
	BatchCancel = "/trade/api/v1/batchCancel"
	//The order information
	GetOrder = "/trade/api/v1/getOrder"
	//Gets an outstanding order
	GetOpenOrders = "/trade/api/v1/getOpenOrders"
	//Get orders history
	GetHistorys = "/trade/api/v1/getHistorys"
	//Gets multiple order information
	GetBatchOrders = "/trade/api/v1/getBatchOrders"
	//Gets the recharge address
	GetPayInAddress = "/trade/api/v1/getPayInAddress"
	//Get the withdrawal address
	GetPayOutAddress = "/trade/api/v1/getPayOutAddress"
	//Gets the recharge record
	GetPayInRecord = "/trade/api/v1/getPayInRecord"
	//Get the withdrawal record
	GetPayOutRecord = "/trade/api/v1/getPayOutRecord"
	//Withdrawal configuration
	GetWithdrawConfig = "/trade/api/v1/getWithdrawConfig"
	//withdraw
	Withdrawal = "/trade/api/v1/withdraw"
)

var KlinePeriodConverter = map[KlinePeriod]string{
	KLINE_PERIOD_1MIN:   "1min",
	KLINE_PERIOD_3MIN:   "3min",
	KLINE_PERIOD_5MIN:   "5min",
	KLINE_PERIOD_15MIN:  "15min",
	KLINE_PERIOD_30MIN:  "30min",
	KLINE_PERIOD_60MIN:  "1hour",
	KLINE_PERIOD_1H:     "1hour",
	KLINE_PERIOD_2H:     "2hour",
	KLINE_PERIOD_4H:     "4hour",
	KLINE_PERIOD_6H:     "6hour",
	KLINE_PERIOD_8H:     "8hour",
	KLINE_PERIOD_12H:    "12hour",
	KLINE_PERIOD_1DAY:   "1day",
	KLINE_PERIOD_3DAY:   "3day",
	KLINE_PERIOD_1WEEK:  "7day",
	KLINE_PERIOD_1MONTH: "30day",
}

type Atop struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

// hao
func (Atop *Atop) buildPostForm(postForm *url.Values) error {
	postForm.Set("accesskey", Atop.accessKey)
	nonce := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	postForm.Set("nonce", nonce)
	payload := postForm.Encode()
	//logs.D("payload", payload)
	sign, _ := GetParamHmacSHA256Sign(Atop.secretKey, payload)
	postForm.Set("signature", sign)
	return nil
}
func New(client *http.Client, apiKey, secretKey string) *Atop {
	return &Atop{apiKey, secretKey, client}
}
func (Atop *Atop) String() string {
	return "atop.com"
}

// hao
func (Atop *Atop) GetTicker(currency CurrencyPair) (*Ticker, error) {
	market := strings.ToLower(currency.String())
	tickerUrl := ApiBaseUrl + fmt.Sprintf(GetTicker, market)
	resp, err := HttpGet(Atop.httpClient, tickerUrl)
	if err != nil {
		return nil, err
	}
	respMap := resp
	var ticker Ticker
	ticker.Pair = currency
	ticker.Date = uint64(time.Now().Unix())
	ticker.Last = num.ToFloat64(respMap["price"])
	ticker.Buy = num.ToFloat64(respMap["bid"])
	ticker.Sell = num.ToFloat64(respMap["ask"])
	ticker.Low = num.ToFloat64(respMap["low"])
	ticker.High = num.ToFloat64(respMap["high"])
	ticker.Vol = num.ToFloat64(respMap["coinVol"])
	return &ticker, nil
}

// hao
func (Atop *Atop) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	market := strings.ToLower(currency.String())
	depthUrl := ApiBaseUrl + fmt.Sprintf(GetDepth, market)
	resp, err := HttpGet(Atop.httpClient, depthUrl)
	if err != nil {
		return nil, err
	}
	respMap := resp
	bids := respMap["bids"].([]any)
	asks := respMap["asks"].([]any)
	depth := new(Depth)
	depth.Pair = currency
	for _, bid := range bids {
		_bid := bid.([]any)
		amount := num.ToFloat64(_bid[1])
		price := num.ToFloat64(_bid[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}
	for _, ask := range asks {
		_ask := ask.([]any)
		amount := num.ToFloat64(_ask[1])
		price := num.ToFloat64(_ask[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}
	sort.Sort(depth.AskList)
	return depth, nil
}

// hao
func (Atop *Atop) plateOrder(amount, price string, currencyPair CurrencyPair, orderType, orderSide string) (*Order, error) {
	pair := Atop.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + PlateOrder
	params := url.Values{}
	params.Set("market", pair.ToLower().String()) //btc_usdt eth_usdt
	if orderSide == "buy" {
		params.Set("type", strconv.Itoa(1))
	} else {
		params.Set("type", strconv.Itoa(0))
	}
	//params.Set("type", orderSide)//Transaction Type  1、buy 0、sell
	params.Set("price", price)
	params.Set("number", amount)
	if orderType == "market" {
		params.Set("entrustType", strconv.Itoa(1))
	} else {
		params.Set("entrustType", strconv.Itoa(0))
	}
	//params.Set("entrustType", orderType)//Delegate type  0、limit，1、market
	Atop.buildPostForm(&params)
	resp, err := HttpPostForm(Atop.httpClient, path, params)
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
	if code != 200 {
		return nil, errors.New(respMap["info"].(string))
	}
	//return &Order{}, nil
	data := respMap["data"].(map[string]any)
	orderId := data["id"].(float64)
	side := BUY
	if orderSide == "sale" {
		side = SELL
	}
	return &Order{
		Currency: pair,
		//OrderID:
		OrderID2:   strconv.FormatFloat(orderId, 'f', 0, 64),
		Price:      num.ToFloat64(price),
		Amount:     num.ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       side,
		Status:     ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}

// hao
func (Atop *Atop) GetAccount() (*Account, error) {
	params := url.Values{}
	Atop.buildPostForm(&params)
	path := ApiBaseUrl + GetBalance
	//logs.D("GetBalance", path)
	resp, err := HttpPostForm(Atop.httpClient, path, params)
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return nil, err
	}
	data := respMap["data"].(map[string]any)
	atc := Account{}
	atc.Exchange = Atop.String()
	atc.SubAccounts = make(map[Currency]SubAccount)
	for k, v := range data {
		cur := NewCurrency(k, "")
		vv := v.(map[string]any)
		sub := SubAccount{}
		sub.Currency = cur
		sub.Amount = num.ToFloat64(vv["available"]) + num.ToFloat64(vv["freeze"])
		sub.ForzenAmount = num.ToFloat64(vv["freeze"])
		atc.SubAccounts[cur] = sub
	}
	return &atc, nil
}

// hao
func (Atop *Atop) LimitBuy(amount, price string, currencyPair CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return Atop.plateOrder(amount, price, currencyPair, "limit", "buy")
}

// hao
func (Atop *Atop) LimitSell(amount, price string, currencyPair CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return Atop.plateOrder(amount, price, currencyPair, "limit", "sale")
}

// hao
func (Atop *Atop) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return Atop.plateOrder(amount, price, currencyPair, "market", "buy")
}

// hao
func (Atop *Atop) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return Atop.plateOrder(amount, price, currencyPair, "market", "sale")
}
func (Atop *Atop) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	currencyPair = Atop.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + CancelOrder
	params := url.Values{}
	params.Set("api_key", Atop.accessKey)
	params.Set("market", currencyPair.ToLower().String())
	params.Set("id", orderId)
	Atop.buildPostForm(&params)
	resp, err := HttpPostForm(Atop.httpClient, path, params)
	if err != nil {
		return false, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		log.Println(string(resp))
		return false, err
	}
	code := respMap["code"].(float64)
	if code != 200 {
		return false, errors.New(respMap["info"].(string))
	}
	//orderIdCanceled := num.ToInt[int](respMap["orderId"])
	//if orderIdCanceled <= 0 {
	//	return false, errors.New(string(resp))
	//}
	return true, nil
}

// hao？
func (Atop *Atop) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	currencyPair = Atop.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + GetOrder
	log.Println(path)
	params := url.Values{}
	params.Set("api_key", Atop.accessKey)
	params.Set("market", currencyPair.ToLower().String())
	params.Set("id", orderId)
	Atop.buildPostForm(&params)
	resp, err := HttpPostForm(Atop.httpClient, path, params)
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
	if code != 200 {
		return nil, errors.New(respMap["info"].(string))
	}
	data := respMap["data"].(map[string]any)
	status := data["status"].(float64)
	side := data["flag"]
	ord := Order{}
	ord.Currency = currencyPair
	//ord.OrderID = num.ToInt[int](orderId)
	ord.OrderID2 = orderId
	if side == "sale" {
		ord.Side = SELL
	} else {
		ord.Side = BUY
	}
	switch status {
	case 0:
		ord.Status = ORDER_UNFINISH
	case 1:
		ord.Status = ORDER_PART_FINISH
	case 2:
		ord.Status = ORDER_FINISH
	case 3:
		ord.Status = ORDER_CANCEL
		//case 4:
		//	ord.Status = new(TradeStatus)//settle
		//case "PENDING_CANCEL":
		//	ord.Status = ORDER_CANCEL_ING
		//case "REJECTED":
		//	ord.Status = ORDER_REJECT
	}
	ord.Amount = num.ToFloat64(data["number"])
	ord.Price = num.ToFloat64(data["price"])
	ord.DealAmount = ord.Amount - num.ToFloat64(data["completeNumber"]) //？
	ord.AvgPrice = num.ToFloat64(data["avg_price"])                     // response no avg price ， fill price
	return &ord, nil
}

// hao
func (Atop *Atop) GetUnfinishedOrders(currencyPair CurrencyPair) ([]Order, error) {
	pair := Atop.adaptCurrencyPair(currencyPair)
	path := ApiBaseUrl + GetOpenOrders
	params := url.Values{}
	params.Set("market", pair.ToLower().String())
	params.Set("page", "1")
	params.Set("pageSize", "10000")
	Atop.buildPostForm(&params)
	resp, err := HttpPostForm(Atop.httpClient, path, params)
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
	if code != 200 {
		return nil, errors.New(respMap["info"].(string))
	}
	data := respMap["data"].([]any)
	orders := make([]Order, 0)
	for _, ord := range data {
		ordData := ord.(map[string]any)
		orderId := strconv.FormatFloat(ordData["id"].(float64), 'f', 0, 64)
		orders = append(orders, Order{
			OrderID:   0,
			OrderID2:  orderId,
			Currency:  currencyPair,
			Price:     num.ToFloat64(ordData["price"]),
			Amount:    num.ToFloat64(ordData["number"]),
			Side:      TradeSide(num.ToInt[int](ordData["type"])),
			Status:    TradeStatus(num.ToInt[int](ordData["status"])),
			OrderTime: num.ToInt[int](ordData["time"])})
	}
	return orders, nil
}

// hao
func (Atop *Atop) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, opt ...OptionalParameter) ([]Kline, error) {
	pair := Atop.adaptCurrencyPair(currency)
	params := url.Values{}
	params.Set("market", pair.ToLower().String())
	//params.Set("type", "1min") //1min,5min,15min,30min,1hour,6hour,1day,7day,30day
	params.Set("type", KlinePeriodConverter[period]) //1min,5min,15min,30min,1hour,6hour,1day,7day,30day
	MergeOptionalParameter(&params, opt...)
	klineUrl := ApiBaseUrl + GetKLine + "?" + params.Encode()
	kLines, err := HttpGet(Atop.httpClient, klineUrl)
	if err != nil {
		return nil, err
	}
	var klineRecords []Kline
	for _, _record := range kLines["datas"].([]any) {
		r := Kline{Pair: currency}
		record := _record.([]any)
		for i, e := range record {
			switch i {
			case 0:
				r.Timestamp = int64(e.(float64)) //to unix timestramp
			case 1:
				r.Open = num.ToFloat64(e)
			case 2:
				r.High = num.ToFloat64(e)
			case 3:
				r.Low = num.ToFloat64(e)
			case 4:
				r.Close = num.ToFloat64(e)
			case 5:
				r.Vol = num.ToFloat64(e)
			}
		}
		klineRecords = append(klineRecords, r)
	}
	return klineRecords, nil
}

// hao Non-individual, transaction record of the entire exchange
func (Atop *Atop) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	pair := Atop.adaptCurrencyPair(currencyPair)
	params := url.Values{}
	params.Set("market", pair.ToLower().String())
	apiUrl := ApiBaseUrl + GetTrades + "?" + params.Encode()
	resp, err := HttpGet(Atop.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}
	var trades []Trade
	for _, v := range resp {
		m := v.(map[string]any)
		ty := SELL
		if m["isBuyerMaker"].(bool) {
			ty = BUY
		}
		trades = append(trades, Trade{
			Tid:    num.ToInt[int64](m["id"]),
			Type:   ty,
			Amount: num.ToFloat64(m["qty"]),
			Price:  num.ToFloat64(m["price"]),
			Date:   num.ToInt[int64](m["time"]),
			Pair:   currencyPair,
		})
	}
	return trades, nil
}
func (Atop *Atop) GetOrderHistorys(currency CurrencyPair, optional ...OptionalParameter) ([]Order, error) {
	//panic("not support")
	pair := Atop.adaptCurrencyPair(currency)
	path := ApiBaseUrl + GetHistorys
	params := url.Values{}
	params.Set("market", pair.ToLower().String())
	MergeOptionalParameter(&params, optional...)
	Atop.buildPostForm(&params)
	resp, err := HttpPostForm(Atop.httpClient, path, params)
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
	if code != 200 {
		return nil, errors.New(respMap["info"].(string))
	}
	data := respMap["data"].(map[string]any)
	records := data["record"].([]any)
	orders := make([]Order, 0)
	for _, ord := range records {
		ordData := ord.(map[string]any)
		orderId := strconv.FormatFloat(ordData["id"].(float64), 'f', 0, 64)
		orders = append(orders, Order{
			OrderID:   0,
			OrderID2:  orderId,
			Currency:  currency,
			Price:     num.ToFloat64(ordData["price"]),
			Amount:    num.ToFloat64(ordData["number"]),
			Side:      TradeSide(num.ToInt[int](ordData["type"])),
			Status:    TradeStatus(num.ToInt[int](ordData["status"])),
			OrderTime: num.ToInt[int](ordData["time"])})
	}
	return orders, nil
}

// hao
func (Atop *Atop) Withdraw(amount, memo string, currency Currency, fees, receiveAddr, safePwd string) (string, error) {
	params := url.Values{}
	coin := strings.ToLower(currency.Symbol)
	path := ApiBaseUrl + Withdrawal
	params.Set("coin", coin)
	params.Set("address", receiveAddr)
	params.Set("amount", amount)
	params.Set("receiveAddr", receiveAddr)
	params.Set("safePwd", safePwd)
	//params.Set("memo", memo)
	Atop.buildPostForm(&params)
	resp, err := HttpPostForm(Atop.httpClient, path, params)
	if err != nil {
		return "", err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return "", err
	}
	if respMap["code"].(float64) == 200 {
		return respMap["id"].(string), nil
	}
	return "", errors.New(string(resp))
}
func (Atop *Atop) CancelWithdraw(id string, currency Currency, safePwd string) (bool, error) {
	panic("not support")
}
func (Atop *Atop) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	return pair
}
