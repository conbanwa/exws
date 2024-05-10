package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/num"
	"github.com/conbanwa/wstrader/util"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type BaseResponse struct {
	Error  []string `json:"error"`
	Result any      `json:"result"`
}
type NewOrderResponse struct {
	Description any      `json:"descr"`
	TxIds       []string `json:"txid"`
}
type Kraken struct {
	httpClient *http.Client
	accessKey,
	secretKey string
}

var (
	BASE_URL   = "https://api.kraken.com"
	API_V0     = "/0/"
	API_DOMAIN = BASE_URL + API_V0
	PUBLIC     = "public/"
	PRIVATE    = "private/"
)

func New(client *http.Client, accesskey, secretkey string) *Kraken {
	return &Kraken{client, accesskey, secretkey}
}
func (k *Kraken) placeOrder(orderType, side, amount, price string, pair CurrencyPair) (*Order, error) {
	apiuri := "private/AddOrder"
	params := url.Values{}
	params.Set("pair", k.convertPair(pair).ToSymbol(""))
	params.Set("type", side)
	params.Set("ordertype", orderType)
	params.Set("price", price)
	params.Set("volume", amount)
	var resp NewOrderResponse
	err := k.doAuthenticatedRequest("POST", apiuri, params, &resp)
	//log.Println
	if err != nil {
		return nil, err
	}
	var tradeSide = SELL
	if "buy" == side {
		tradeSide = BUY
	}
	return &Order{
		Currency: pair,
		OrderID2: resp.TxIds[0],
		Amount:   num.ToFloat64(amount),
		Price:    num.ToFloat64(price),
		Side:     tradeSide,
		Status:   ORDER_UNFINISH}, nil
}
func (k *Kraken) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return k.placeOrder("limit", "buy", amount, price, currency)
}
func (k *Kraken) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	return k.placeOrder("limit", "sell", amount, price, currency)
}
func (k *Kraken) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	return k.placeOrder("market", "buy", amount, price, currency)
}
func (k *Kraken) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	return k.placeOrder("market", "sell", amount, price, currency)
}
func (k *Kraken) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	params := url.Values{}
	apiuri := "private/CancelOrder"
	params.Set("txid", orderId)
	var respMap map[string]any
	err := k.doAuthenticatedRequest("POST", apiuri, params, &respMap)
	if err != nil {
		return false, err
	}
	//log.Println(respMap)
	return true, nil
}
func (k *Kraken) toOrder(orderinfo any) Order {
	omap := orderinfo.(map[string]any)
	descmap := omap["descr"].(map[string]any)
	return Order{
		Amount:     num.ToFloat64(omap["vol"]),
		Price:      num.ToFloat64(descmap["price"]),
		DealAmount: num.ToFloat64(omap["vol_exec"]),
		AvgPrice:   num.ToFloat64(omap["price"]),
		Side:       util.AdaptTradeSide(descmap["type"].(string)),
		Status:     k.convertOrderStatus(omap["status"].(string)),
		Fee:        num.ToFloat64(omap["fee"]),
		OrderTime:  num.ToInt[int](omap["opentm"]),
	}
}
func (k *Kraken) GetOrderInfos(txids ...string) ([]Order, error) {
	params := url.Values{}
	params.Set("txid", strings.Join(txids, ","))
	var resultmap map[string]any
	err := k.doAuthenticatedRequest("POST", "private/QueryOrders", params, &resultmap)
	if err != nil {
		return nil, err
	}
	//log.Println(resultmap)
	var ords []Order
	for txid, v := range resultmap {
		ord := k.toOrder(v)
		ord.OrderID2 = txid
		ords = append(ords, ord)
	}
	return ords, nil
}
func (k *Kraken) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	orders, err := k.GetOrderInfos(orderId)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("not fund the order " + orderId)
	}
	ord := &orders[0]
	ord.Currency = currency
	return ord, nil
}
func (k *Kraken) GetUnfinishedOrders(currency CurrencyPair) ([]Order, error) {
	var result struct {
		Open map[string]any `json:"open"`
	}
	err := k.doAuthenticatedRequest("POST", "private/OpenOrders", url.Values{}, &result)
	if err != nil {
		return nil, err
	}
	var orders []Order
	for txid, v := range result.Open {
		ord := k.toOrder(v)
		ord.OrderID2 = txid
		ord.Currency = currency
		orders = append(orders, ord)
	}
	return orders, nil
}
func (k *Kraken) GetOrderHistorys(currency CurrencyPair, optional ...OptionalParameter) ([]Order, error) {
	panic("")
}
func (k *Kraken) GetAccount() (*Account, error) {
	params := url.Values{}
	apiuri := "private/Balance"
	var resustmap map[string]any
	err := k.doAuthenticatedRequest("POST", apiuri, params, &resustmap)
	if err != nil {
		return nil, err
	}
	acc := new(Account)
	acc.Exchange = k.String()
	acc.SubAccounts = make(map[Currency]SubAccount)
	for key, v := range resustmap {
		currency := k.convertCurrency(key)
		amount := num.ToFloat64(v)
		//log.Println(symbol, amount)
		acc.SubAccounts[currency] = SubAccount{Currency: currency, Amount: amount, ForzenAmount: 0, LoanAmount: 0}
		if currency.Symbol == "XBT" { // adapt to btc
			acc.SubAccounts[BTC] = SubAccount{Currency: BTC, Amount: amount, ForzenAmount: 0, LoanAmount: 0}
		}
	}
	return acc, nil
}

//	func (k *Kraken) GetTradeBalance() {
//		var resultmap map[string]any
//		k.doAuthenticatedRequest("POST", "private/TradeBalance", url.Values{}, &resultmap)
//		log.Println(resultmap)
//	}
func (k *Kraken) GetTicker(currency CurrencyPair) (*Ticker, error) {
	var resultmap map[string]any
	err := k.doAuthenticatedRequest("GET", "public/Ticker?pair="+k.convertPair(currency).ToSymbol(""), url.Values{}, &resultmap)
	if err != nil {
		return nil, err
	}
	ticker := new(Ticker)
	ticker.Pair = currency
	for _, t := range resultmap {
		tickermap := t.(map[string]any)
		ticker.Last = num.ToFloat64(tickermap["c"].([]any)[0])
		ticker.Buy = num.ToFloat64(tickermap["b"].([]any)[0])
		ticker.Sell = num.ToFloat64(tickermap["a"].([]any)[0])
		ticker.Low = num.ToFloat64(tickermap["l"].([]any)[0])
		ticker.High = num.ToFloat64(tickermap["h"].([]any)[0])
		ticker.Vol = num.ToFloat64(tickermap["v"].([]any)[0])
	}
	return ticker, nil
}
func (k *Kraken) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	apiuri := fmt.Sprintf("public/Depth?pair=%s&count=%d", k.convertPair(currency).ToSymbol(""), size)
	var resultmap map[string]any
	err := k.doAuthenticatedRequest("GET", apiuri, url.Values{}, &resultmap)
	if err != nil {
		return nil, err
	}
	//log.Println(respMap)
	dep := Depth{}
	dep.Pair = currency
	for _, d := range resultmap {
		depmap := d.(map[string]any)
		asksmap := depmap["asks"].([]any)
		bidsmap := depmap["bids"].([]any)
		for _, v := range asksmap {
			ask := v.([]any)
			dep.AskList = append(dep.AskList, DepthRecord{Price: num.ToFloat64(ask[0]), Amount: num.ToFloat64(ask[1])})
		}
		for _, v := range bidsmap {
			bid := v.([]any)
			dep.BidList = append(dep.BidList, DepthRecord{Price: num.ToFloat64(bid[0]), Amount: num.ToFloat64(bid[1])})
		}
		break
	}
	sort.Sort(sort.Reverse(dep.AskList)) //reverse
	return &dep, nil
}
func (k *Kraken) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, opt ...OptionalParameter) ([]Kline, error) {
	panic("")
}

// 非个人，整个交易所的交易记录
func (k *Kraken) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("")
}
func (k *Kraken) String() string {
	return KRAKEN
}
func (k *Kraken) buildParamsSigned(apiuri string, postForm *url.Values) string {
	postForm.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))
	urlPath := API_V0 + apiuri
	secretByte, _ := base64.StdEncoding.DecodeString(k.secretKey)
	encode := []byte(postForm.Get("nonce") + postForm.Encode())
	sha := sha256.New()
	sha.Write(encode)
	shaSum := sha.Sum(nil)
	pathSha := append([]byte(urlPath), shaSum...)
	mac := hmac.New(sha512.New, secretByte)
	mac.Write(pathSha)
	macSum := mac.Sum(nil)
	sign := base64.StdEncoding.EncodeToString(macSum)
	return sign
}
func (k *Kraken) doAuthenticatedRequest(method, apiuri string, params url.Values, ret any) error {
	headers := map[string]string{}
	if "POST" == method {
		signature := k.buildParamsSigned(apiuri, &params)
		headers = map[string]string{
			"API-Key":  k.accessKey,
			"API-Sign": signature,
		}
	}
	resp, err := NewRequest(k.httpClient, method, API_DOMAIN+apiuri, params.Encode(), headers)
	if err != nil {
		return err
	}
	//println(string(resp))
	var base BaseResponse
	base.Result = ret
	err = json.Unmarshal(resp, &base)
	if err != nil {
		return err
	}
	//println(string(resp))
	if len(base.Error) > 0 {
		return errors.New(base.Error[0])
	}
	return nil
}
func (k *Kraken) convertCurrency(currencySymbol string) Currency {
	if len(currencySymbol) >= 4 {
		currencySymbol = strings.Replace(currencySymbol, "X", "", 1)
		currencySymbol = strings.Replace(currencySymbol, "Z", "", 1)
	}
	return NewCurrency(currencySymbol, "")
}
func (k *Kraken) convertPair(pair CurrencyPair) CurrencyPair {
	if "BTC" == pair.CurrencyA.Symbol {
		return NewCurrencyPair(XBT, pair.CurrencyB)
	}
	if "BTC" == pair.CurrencyB.Symbol {
		return NewCurrencyPair(pair.CurrencyA, XBT)
	}
	return pair
}
func (k *Kraken) convertOrderStatus(status string) TradeStatus {
	switch status {
	case "open", "pending":
		return ORDER_UNFINISH
	case "canceled", "expired":
		return ORDER_CANCEL
	case "filled", "closed":
		return ORDER_FINISH
	case "partialfilled":
		return ORDER_PART_FINISH
	}
	return ORDER_UNFINISH
}
