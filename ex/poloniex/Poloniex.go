package poloniex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/conbanwa/logs"
)

const EXCHANGE_NAME = "poloniex.com"
const (
	BASE_URL       = "https://poloniex.com/"
	TRADE_API      = BASE_URL + "tradingApi"
	PUBLIC_URL     = BASE_URL + "public"
	TICKER_API     = "?command=returnTicker"
	CURRENCIES_API = "?command=returnCurrencies"
	ORDER_BOOK_API = "?command=returnOrderBook&currencyPair=%s&depth=%d"
)

type Poloniex struct {
	accessKey,
	secretKey string
	client *http.Client
}

func New(client *http.Client, accessKey, secretKey string) *Poloniex {
	return &Poloniex{accessKey, secretKey, client}
}
func (poloniex *Poloniex) String() string {
	return POLONIEX
}
func (poloniex *Poloniex) GetTicker(currency CurrencyPair) (*Ticker, error) {
	//logs.E(poloniex.adaptCurrencyPair(currency).ToSymbol2("_"))
	respMap, err := HttpGet(poloniex.client, PUBLIC_URL+TICKER_API)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	pair := currency.AdaptUsdToUsdt().Reverse().ToSymbol("_")
	//println(pair)
	tickermap, ok := respMap[pair].(map[string]any)
	if !ok {
		return new(Ticker), errors.New("not found")
	}
	ticker := new(Ticker)
	ticker.Pair = currency
	ticker.High, _ = strconv.ParseFloat(tickermap["high24hr"].(string), 64)
	ticker.Low, _ = strconv.ParseFloat(tickermap["low24hr"].(string), 64)
	ticker.Last, _ = strconv.ParseFloat(tickermap["last"].(string), 64)
	ticker.Buy, _ = strconv.ParseFloat(tickermap["highestBid"].(string), 64)
	ticker.Sell, _ = strconv.ParseFloat(tickermap["lowestAsk"].(string), 64)
	ticker.Vol, _ = strconv.ParseFloat(tickermap["quoteVolume"].(string), 64)
	//logs.E(tickermap)
	return ticker, nil
}
func (poloniex *Poloniex) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	respMap, err := HttpGet(poloniex.client, PUBLIC_URL+
		fmt.Sprintf(ORDER_BOOK_API, currency.AdaptUsdToUsdt().Reverse().ToSymbol("_"), size))
	if err != nil {
		logs.E(err)
		return nil, err
	}
	if respMap["asks"] == nil {
		logs.E(respMap)
		return nil, errors.New(fmt.Sprintf("%+v", respMap))
	}
	if _, isOK := respMap["asks"].([]any); !isOK {
		logs.E(respMap)
		return nil, errors.New(fmt.Sprintf("%+v", respMap))
	}
	var depth Depth
	depth.Pair = currency
	for _, v := range respMap["asks"].([]any) {
		var dr DepthRecord
		for i, vv := range v.([]any) {
			switch i {
			case 0:
				dr.Price, _ = strconv.ParseFloat(vv.(string), 64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.AskList = append(depth.AskList, dr)
	}
	for _, v := range respMap["bids"].([]any) {
		var dr DepthRecord
		for i, vv := range v.([]any) {
			switch i {
			case 0:
				dr.Price, _ = strconv.ParseFloat(vv.(string), 64)
			case 1:
				dr.Amount = vv.(float64)
			}
		}
		depth.BidList = append(depth.BidList, dr)
	}
	return &depth, nil
}
func (poloniex *Poloniex) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]Kline, error) {
	return nil, nil
}
func (poloniex *Poloniex) placeLimitOrder(command, amount, price string, currency CurrencyPair) (*Order, error) {
	postData := url.Values{}
	postData.Set("command", command)
	postData.Set("currencyPair", currency.AdaptUsdToUsdt().Reverse().ToSymbol("_"))
	postData.Set("rate", price)
	postData.Set("amount", amount)
	sign, _ := poloniex.buildPostForm(&postData)
	headers := map[string]string{
		"Key":  poloniex.accessKey,
		"Sign": sign}
	resp, err := HttpPostForm2(poloniex.client, TRADE_API, postData, headers)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil || respMap["error"] != nil {
		logs.E(err, string(resp))
		return nil, err
	}
	orderNumber := respMap["orderNumber"].(string)
	order := new(Order)
	order.OrderTime = int(time.Now().Unix() * 1000)
	order.OrderID, _ = strconv.Atoi(orderNumber)
	order.OrderID2 = orderNumber
	order.Amount, _ = strconv.ParseFloat(amount, 64)
	order.Price, _ = strconv.ParseFloat(price, 64)
	order.Status = ORDER_UNFINISH
	order.Currency = currency
	switch command {
	case "sell":
		order.Side = SELL
	case "buy":
		order.Side = BUY
	}
	//logs.E(string(resp))
	return order, nil
}
func (poloniex *Poloniex) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return poloniex.placeLimitOrder("buy", amount, price, currency)
}
func (poloniex *Poloniex) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return poloniex.placeLimitOrder("sell", amount, price, currency)
}
func (poloniex *Poloniex) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	postData := url.Values{}
	postData.Set("command", "cancelOrder")
	postData.Set("orderNumber", orderId)
	sign, err := poloniex.buildPostForm(&postData)
	if err != nil {
		logs.E(err)
		return false, err
	}
	headers := map[string]string{
		"Key":  poloniex.accessKey,
		"Sign": sign}
	resp, err := HttpPostForm2(poloniex.client, TRADE_API, postData, headers)
	if err != nil {
		logs.E(err)
		return false, err
	}
	//logs.E(string(resp));
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil || respMap["error"] != nil {
		//logs.E(err, string(resp))
		return false, errors.New(string(resp))
	}
	if success := int(respMap["success"].(float64)); success != 1 {
		logs.E(respMap)
		return false, nil
	}
	return true, nil
}
func (poloniex *Poloniex) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	postData := url.Values{}
	postData.Set("command", "returnOrderTrades")
	postData.Set("orderNumber", orderId)
	sign, _ := poloniex.buildPostForm(&postData)
	headers := map[string]string{
		"Key":  poloniex.accessKey,
		"Sign": sign}
	resp, err := HttpPostForm2(poloniex.client, TRADE_API, postData, headers)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	//println(string(resp))
	if strings.Contains(string(resp), "error") {
		if ords, err1 := poloniex.GetUnfinishedOrders(currency); err1 != nil {
			logs.E(err1)
			if strings.Contains(err1.Error(), "Order not found") {
				return nil, EX_ERR_NOT_FIND_ORDER
			}
		} else {
			_ordId, _ := strconv.Atoi(orderId)
			for _, ord := range ords {
				if ord.OrderID == _ordId {
					return &ord, nil
				}
			}
		}
		//logs.E(string(resp))
		return nil, errors.New(string(resp))
	}
	respMap := make([]any, 0)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		logs.E(err, string(resp))
		return nil, err
	}
	order := new(Order)
	order.OrderID, _ = strconv.Atoi(orderId)
	order.Currency = currency
	total := 0.0
	for _, v := range respMap {
		vv := v.(map[string]any)
		_amount, _ := strconv.ParseFloat(vv["amount"].(string), 64)
		_rate, _ := strconv.ParseFloat(vv["rate"].(string), 64)
		_fee, _ := strconv.ParseFloat(vv["fee"].(string), 64)
		order.DealAmount += _amount
		total += _amount * _rate
		order.Fee = _fee
		if strings.Compare("sell", vv["type"].(string)) == 0 {
			order.Side = SELL
		} else {
			order.Side = BUY
		}
	}
	order.AvgPrice = total / order.DealAmount
	return order, nil
}
func (poloniex *Poloniex) GetUnfinishedOrders(currency CurrencyPair) ([]Order, error) {
	postData := url.Values{}
	postData.Set("command", "returnOpenOrders")
	postData.Set("currencyPair", currency.AdaptUsdToUsdt().Reverse().ToSymbol("_"))
	sign, err := poloniex.buildPostForm(&postData)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	headers := map[string]string{
		"Key":  poloniex.accessKey,
		"Sign": sign}
	resp, err := HttpPostForm2(poloniex.client, TRADE_API, postData, headers)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	orderAr := make([]any, 1)
	err = json.Unmarshal(resp, &orderAr)
	if err != nil {
		logs.E(err, string(resp))
		return nil, err
	}
	orders := make([]Order, 0)
	for _, v := range orderAr {
		vv := v.(map[string]any)
		order := Order{}
		order.Currency = currency
		order.OrderID, _ = strconv.Atoi(vv["orderNumber"].(string))
		order.OrderID2 = vv["orderNumber"].(string)
		order.Amount, _ = strconv.ParseFloat(vv["amount"].(string), 64)
		order.Price, _ = strconv.ParseFloat(vv["rate"].(string), 64)
		order.Status = ORDER_UNFINISH
		side := vv["type"].(string)
		switch side {
		case "buy":
			order.Side = BUY
		case "sell":
			order.Side = SELL
		}
		orders = append(orders, order)
	}
	//logs.E(orders)
	return orders, nil
}
func (poloniex *Poloniex) GetOrderHistorys(currency CurrencyPair, opt ...OptionalParameter) ([]Order, error) {
	return nil, nil
}
func (poloniex *Poloniex) GetAccount() (*Account, error) {
	postData := url.Values{}
	postData.Add("command", "returnCompleteBalances")
	sign, err := poloniex.buildPostForm(&postData)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Key":  poloniex.accessKey,
		"Sign": sign}
	resp, err := HttpPostForm2(poloniex.client, TRADE_API, postData, headers)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil || respMap["error"] != nil {
		logs.E(err)
		return nil, err
	}
	acc := new(Account)
	acc.Exchange = EXCHANGE_NAME
	acc.SubAccounts = make(map[Currency]SubAccount)
	for k, v := range respMap {
		var currency = NewCurrency(k, "")
		vv := v.(map[string]any)
		subAcc := SubAccount{}
		subAcc.Currency = currency
		subAcc.Amount, _ = strconv.ParseFloat(vv["available"].(string), 64)
		subAcc.ForzenAmount, _ = strconv.ParseFloat(vv["onOrders"].(string), 64)
		acc.SubAccounts[subAcc.Currency] = subAcc
		if currency.Symbol == "USDT" {
			acc.SubAccounts[USD] = subAcc
		}
	}
	return acc, nil
}
func (poloniex *Poloniex) Withdraw(amount string, currency Currency, fees, receiveAddr, safePwd string) (string, error) {
	if currency == BCC {
		currency = BCH
	}
	params := url.Values{}
	params.Add("command", "withdraw")
	params.Add("address", receiveAddr)
	params.Add("amount", amount)
	params.Add("currency", strings.ToUpper(currency.String()))
	sign, err := poloniex.buildPostForm(&params)
	if err != nil {
		return "", err
	}
	headers := map[string]string{
		"Key":  poloniex.accessKey,
		"Sign": sign}
	resp, err := HttpPostForm2(poloniex.client, TRADE_API, params, headers)
	if err != nil {
		logs.E(err)
		return "", err
	}
	println(string(resp))
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		logs.E(err)
		return "", err
	}
	if respMap["error"] == nil {
		return string(resp), nil
	}
	return "", errors.New(string(resp))
}

type DepositsWithdrawals struct {
	Deposits []struct {
		Currency      string    `json:"currency"`
		Address       string    `json:"address"`
		Amount        float64   `json:"amount,string"`
		Confirmations int       `json:"confirmations"`
		TransactionID string    `json:"txid"`
		Timestamp     time.Time `json:"timestamp"`
		Status        string    `json:"status"`
	} `json:"deposits"`
	Withdrawals []struct {
		WithdrawalNumber int64     `json:"withdrawalNumber"`
		Currency         string    `json:"currency"`
		Address          string    `json:"address"`
		Amount           float64   `json:"amount,string"`
		Confirmations    int       `json:"confirmations"`
		TransactionID    string    `json:"txid"`
		Timestamp        time.Time `json:"timestamp"`
		Status           string    `json:"status"`
		IPAddress        string    `json:"ipAddress"`
	} `json:"withdrawals"`
}

func (poloniex *Poloniex) GetDepositsWithdrawals(start, end string) (*DepositsWithdrawals, error) {
	params := url.Values{}
	params.Set("command", "returnDepositsWithdrawals")
	println(start)
	if start != "" {
		params.Set("start", start)
	} else {
		params.Set("start", "0")
	}
	if end != "" {
		params.Set("end", end)
	} else {
		params.Set("end", strconv.FormatInt(time.Now().Unix(), 10))
	}
	sign, err := poloniex.buildPostForm(&params)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Key":  poloniex.accessKey,
		"Sign": sign}
	resp, err := HttpPostForm2(poloniex.client, TRADE_API, params, headers)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	println(string(resp))
	records := new(DepositsWithdrawals)
	err = json.Unmarshal(resp, records)
	return records, err
}
func (poloniex *Poloniex) buildPostForm(postForm *url.Values) (string, error) {
	postForm.Add("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	payload := postForm.Encode()
	//println(payload)
	sign, err := GetParamHmacSHA512Sign(poloniex.secretKey, payload)
	if err != nil {
		return "", err
	}
	//logs.E(sign)
	return sign, nil
}
func (poloniex *Poloniex) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("unimplements")
}
func (poloniex *Poloniex) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unsupport the market order")
}
func (poloniex *Poloniex) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("unsupport the market order")
}
func (poloniex *Poloniex) GetAllCurrencies() (map[string]*PoloniexCurrency, error) {
	respMap, err := HttpGet(poloniex.client, PUBLIC_URL+CURRENCIES_API)
	if err != nil || respMap["error"] != nil {
		logs.E(err)
		return nil, err
	}
	result := map[string]*PoloniexCurrency{}
	for k, v := range respMap {
		currencyMap := v.(map[string]any)
		poloniexCurrency := new(PoloniexCurrency)
		poloniexCurrency.ID = int(currencyMap["id"].(float64))
		poloniexCurrency.Name, _ = currencyMap["name"].(string)
		poloniexCurrency.TxFee, _ = currencyMap["txFee"].(string)
		poloniexCurrency.MinConf = int(currencyMap["minConf"].(float64))
		poloniexCurrency.DepositAddress, _ = currencyMap["depositAddress"].(string)
		poloniexCurrency.Disabled = int(currencyMap["disabled"].(float64))
		poloniexCurrency.Delisted = int(currencyMap["delisted"].(float64))
		poloniexCurrency.Frozen = int(currencyMap["frozen"].(float64))
		poloniexCurrency.HumanType, _ = currencyMap["humanType"].(string)
		poloniexCurrency.CurrencyType, _ = currencyMap["currencyType"].(string)
		poloniexCurrency.Blockchain, _ = currencyMap["blockchain"].(string)
		poloniexCurrency.IsGeofenced = int(currencyMap["isGeofenced"].(float64))
		result[k] = poloniexCurrency
	}
	return result, nil
}
