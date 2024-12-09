package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/num"
	"github.com/conbanwa/slice"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"github.com/conbanwa/wstrader/util"
	"github.com/conbanwa/wstrader/web"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	GlobalApiBaseUrl           = "https://api.binance.com"
	UsApiBaseUrl               = "https://api.binance.us"
	JeApiBaseUrl               = "https://api.binance.je"
	ApiV3                      = GlobalApiBaseUrl + "/api/v3/"
	TickerUri                  = "ticker/24hr?symbol=%s"
	TickersUri                 = "ticker/allBookTickers"
	FeeUrl                     = GlobalApiBaseUrl + "/sapi/v1/asset/tradeFee"
	ExWithdrawUrl              = GlobalApiBaseUrl + "/sapi/v1/asset/assetDetail"
	NetworkWithdrawUrl         = GlobalApiBaseUrl + "/sapi/v1/capital/config/getall"
	DepthUri                   = "depth?symbol=%s&limit=%d"
	AccountUri                 = "account?"
	OrderUri                   = "order"
	UnfinishedOrdersInfo       = "openOrders?"
	KlineUri                   = "klines"
	ServerTimeUrl              = "time"
	FutureUsdWsBaseUrl         = "wss://fstream.binance.com/ws"
	FutureCoinWsBaseUrl        = "wss://dstream.binance.com/ws"
	TestnetSpotApiBaseUrl      = GlobalApiBaseUrl
	TestnetSpotWsBaseUrl       = "wss://testnet.binance.vision/ws"
	TestnetSpotStreamBaseUrl   = "wss://testnet.binance.vision/stream"
	TestnetFutureUsdBaseUrl    = "https://testnet.binancefuture.com"
	TestnetFutureUsdWsBaseUrl  = "wss://fstream.binance.com/ws"
	TestnetFutureCoinWsBaseUrl = "wss://dstream.binance.com/ws"
)

var internalKlinePeriodConverter = map[KlinePeriod]string{
	KLINE_PERIOD_1MIN:   "1m",
	KLINE_PERIOD_3MIN:   "3m",
	KLINE_PERIOD_5MIN:   "5m",
	KLINE_PERIOD_15MIN:  "15m",
	KLINE_PERIOD_30MIN:  "30m",
	KLINE_PERIOD_60MIN:  "1h",
	KLINE_PERIOD_1H:     "1h",
	KLINE_PERIOD_2H:     "2h",
	KLINE_PERIOD_4H:     "4h",
	KLINE_PERIOD_6H:     "6h",
	KLINE_PERIOD_8H:     "8h",
	KLINE_PERIOD_12H:    "12h",
	KLINE_PERIOD_1DAY:   "1d",
	KLINE_PERIOD_3DAY:   "3d",
	KLINE_PERIOD_1WEEK:  "1w",
	KLINE_PERIOD_1MONTH: "1M",
}

type Filter struct {
	FilterType          string  `json:"filterType"`
	MaxPrice            float64 `json:"maxPrice,string"`
	MinPrice            float64 `json:"minPrice,string"`
	TickSize            float64 `json:"tickSize,string"`
	MultiplierUp        float64 `json:"multiplierUp,string"`
	MultiplierDown      float64 `json:"multiplierDown,string"`
	AvgPriceMins        int     `json:"avgPriceMins"`
	MinQty              float64 `json:"minQty,string"`
	MaxQty              float64 `json:"maxQty,string"`
	StepSize            float64 `json:"stepSize,string"`
	MinNotional         float64 `json:"minNotional,string"`
	ApplyToMarket       bool    `json:"applyToMarket"`
	Limit               int     `json:"limit"`
	MaxNumAlgoOrders    int     `json:"maxNumAlgoOrders"`
	MaxNumIcebergOrders int     `json:"maxNumIcebergOrders"`
	MaxNumOrders        int     `json:"maxNumOrders"`
}
type RateLimit struct {
	Interval      string `json:"interval"`
	IntervalNum   int64  `json:"intervalNum"`
	Limit         int64  `json:"limit"`
	RateLimitType string `json:"rateLimitType"`
}
type TradeSymbol struct {
	Symbol                     string   `json:"symbol"`
	Status                     string   `json:"status"`
	BaseAsset                  string   `json:"baseAsset"`
	BaseAssetPrecision         int      `json:"baseAssetPrecision"`
	QuoteAsset                 string   `json:"quoteAsset"`
	QuotePrecision             int      `json:"quotePrecision"`
	QuoteAssetPrecision        int      `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int      `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int      `json:"quoteCommissionPrecision"`
	Filters                    []Filter `json:"filters"`
	IcebergAllowed             bool     `json:"icebergAllowed"`
	IsMarginTradingAllowed     bool     `json:"isMarginTradingAllowed"`
	IsSpotTradingAllowed       bool     `json:"isSpotTradingAllowed"`
	OcoAllowed                 bool     `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool     `json:"quoteOrderQtyMarketAllowed"`
	OrderTypes                 []string `json:"orderTypes"`
}

func (ts TradeSymbol) GetMinBase() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "LOT_SIZE" {
			return v.MinQty
		}
	}
	return 0
}
func (ts TradeSymbol) GetBaseStepN() int {
	for _, v := range ts.Filters {
		if v.FilterType == "LOT_SIZE" {
			step := strconv.FormatFloat(v.StepSize, 'f', -1, 64)
			pres := strings.Split(step, ".")
			if len(pres) == 1 {
				return 0
			}
			return len(pres[1])
		}
	}
	return 8
}
func (ts TradeSymbol) GetBaseStep() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "LOT_SIZE" {
			return v.StepSize
		}
	}
	panic(ts.Filters)
	return 1
}
func (ts TradeSymbol) GetMinQuote() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "NOTIONAL" {
			return v.MinNotional
		}
		if v.FilterType == "MIN_NOTIONAL" {
			panic(v.MinNotional)
		}
	}
	return 0
}
func (ts TradeSymbol) GetPriceStep() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "PRICE_FILTER" {
			return v.TickSize
			//step := strconv.FormatFloat(v.TickSize, 'f', -1, 64)
			//pres := strings.Split(step, ".")
			//if len(pres) == 1 {
			//	return 0
			//}
			//return len(pres[1])
		}
	}
	return 0
}
func (ts TradeSymbol) GetMinPrice() float64 {
	for _, v := range ts.Filters {
		if v.FilterType == "PRICE_FILTER" {
			return v.MinPrice
		}
	}
	return 0
}

type ExchangeInfo struct {
	Timezone        string        `json:"timezone"`
	ServerTime      int           `json:"serverTime"`
	ExchangeFilters []any         `json:"exchangeFilters,omitempty"`
	RateLimits      []RateLimit   `json:"rateLimits"`
	Symbols         []TradeSymbol `json:"symbols"`
}
type Binance struct {
	accessKey  string
	secretKey  string
	baseUrl    string
	apiV1      string
	apiV3      string
	httpClient *http.Client
	timeOffset int64 //nanosecond
	*ExchangeInfo
}

func (bn *Binance) buildParamsSigned(postForm *url.Values) (err error) {
	postForm.Set("recvWindow", "60000")
	tonce := strconv.FormatInt(time.Now().UnixNano()+bn.timeOffset, 10)[0:13]
	postForm.Set("timestamp", tonce)
	payload := postForm.Encode()
	sign, err := GetParamHmacSHA256Sign(bn.secretKey, payload)
	postForm.Set("signature", sign)
	return
}
func New(client *http.Client, apiKey, secretKey string) *Binance {
	return NewWithConfig(&APIConfig{
		HttpClient:   client,
		Endpoint:     GlobalApiBaseUrl,
		ApiKey:       apiKey,
		ApiSecretKey: secretKey})
}
func NewWithConfig(config *APIConfig) *Binance {
	if config.Endpoint == "" {
		config.Endpoint = GlobalApiBaseUrl
	}
	bn := &Binance{
		baseUrl:    config.Endpoint,
		apiV1:      config.Endpoint + "/api/v1/",
		apiV3:      config.Endpoint + "/api/v3/",
		accessKey:  config.ApiKey,
		secretKey:  config.ApiSecretKey,
		httpClient: config.HttpClient}
	bn.setTimeOffset()
	return bn
}
func (bn *Binance) String() string {
	return BINANCE
}
func (bn *Binance) Ping() bool {
	if _, err := web.HttpGet(bn.httpClient, bn.apiV3+"ping"); err != nil {
		return false
	}
	return true
}
func (bn *Binance) setTimeOffset() error {
	respMap, err := web.HttpGet(bn.httpClient, bn.apiV3+ServerTimeUrl)
	if err != nil {
		return err
	}
	stime := int64(num.ToInt[int](respMap["serverTime"]))
	st := time.Unix(stime/1000, 1000000*(stime%1000))
	lt := time.Now()
	offset := st.Sub(lt).Nanoseconds()
	bn.timeOffset = offset
	return nil
}
func (bn *Binance) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := bn.apiV3 + fmt.Sprintf(TickerUri, currency.ToSymbol(""))
	tickerMap, err := web.HttpGet(bn.httpClient, tickerUri)
	if err != nil {
		return nil, err
	}
	var ticker Ticker
	ticker.Pair = currency
	t, _ := tickerMap["closeTime"].(float64)
	ticker.Date = uint64(t / 1000)
	ticker.Last = num.ToFloat64(tickerMap["lastPrice"])
	ticker.Buy = num.ToFloat64(tickerMap["bidPrice"])
	ticker.Sell = num.ToFloat64(tickerMap["askPrice"])
	ticker.Low = num.ToFloat64(tickerMap["lowPrice"])
	ticker.High = num.ToFloat64(tickerMap["highPrice"])
	ticker.Vol = num.ToFloat64(tickerMap["volume"])
	return &ticker, nil
}
func (bn *Binance) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	if size <= 5 {
		size = 5
	} else if size <= 10 {
		size = 10
	} else if size <= 20 {
		size = 20
	} else if size <= 50 {
		size = 50
	} else if size <= 100 {
		size = 100
	} else if size <= 500 {
		size = 500
	} else {
		size = 1000
	}
	apiUrl := fmt.Sprintf(bn.apiV3+DepthUri, currencyPair.ToSymbol(""), size)
	resp, err := web.HttpGet(bn.httpClient, apiUrl)
	if err != nil {
		return nil, err
	}
	if _, isok := resp["code"]; isok {
		return nil, errors.New(resp["msg"].(string))
	}
	bids := resp["bids"].([]any)
	asks := resp["asks"].([]any)
	depth := new(Depth)
	depth.Pair = currencyPair
	depth.UTime = time.Now()
	n := 0
	for _, bid := range bids {
		_bid := bid.([]any)
		amount := num.ToFloat64(_bid[1])
		price := num.ToFloat64(_bid[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
		n++
		if n == size {
			break
		}
	}
	n = 0
	for _, ask := range asks {
		_ask := ask.([]any)
		amount := num.ToFloat64(_ask[1])
		price := num.ToFloat64(_ask[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
		n++
		if n == size {
			break
		}
	}
	sort.Sort(sort.Reverse(depth.AskList))
	return depth, nil
}
func (bn *Binance) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (*q.Order, error) {
	// logs.W(amount, price, pair, orderType, orderSide)
	path := bn.apiV3 + OrderUri
	params := url.Values{}
	params.Set("symbol", pair.ToSymbol(""))
	params.Set("side", orderSide)
	params.Set("type", orderType)
	params.Set("newOrderRespType", "ACK")
	params.Set("quantity", amount)
	switch orderType {
	case "LIMIT":
		params.Set("timeInForce", "GTC")
		params.Set("price", price)
	case "MARKET":
		params.Set("newOrderRespType", "RESULT")
	}
	bn.buildParamsSigned(&params)
	resp, err := web.HttpPostForm2(bn.httpClient, path, params, bn.header())
	if err != nil {
		return nil, err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return nil, err
	}
	orderId := num.ToInt[int](respMap["orderId"])
	if orderId <= 0 {
		return nil, errors.New(slice.Bytes2String(resp))
	}
	side := BUY
	if orderSide == "SELL" {
		side = SELL
	}
	dealAmount := num.ToFloat64(respMap["executedQty"])
	cummulativeQuoteQty := num.ToFloat64(respMap["cummulativeQuoteQty"])
	avgPrice := 0.0
	if cummulativeQuoteQty > 0 && dealAmount > 0 {
		avgPrice = cummulativeQuoteQty / dealAmount
	}
	// log.Println(dealAmount, avgPrice)
	return &q.Order{
		Currency:   pair,
		OrderID:    orderId,
		OrderID2:   strconv.Itoa(orderId),
		Price:      num.ToFloat64(price),
		Amount:     num.ToFloat64(amount),
		DealAmount: dealAmount,
		AvgPrice:   avgPrice,
		Side:       side,
		Status:     ORDER_UNFINISH,
		OrderTime:  num.ToInt[int](respMap["transactTime"])}, nil
}
func (bn *Binance) GetAccount() (*Account, error) {
	params := url.Values{}
	bn.buildParamsSigned(&params)
	path := bn.apiV3 + AccountUri + params.Encode()
	respMap, err := web.HttpGet2(bn.httpClient, path, bn.header())
	if err != nil {
		return nil, err
	}
	if _, isok := respMap["code"]; isok {
		return nil, errors.New(respMap["msg"].(string))
	}
	acc := Account{}
	acc.Exchange = bn.String()
	acc.SubAccounts = make(map[Currency]SubAccount)
	balances := respMap["balances"].([]any)
	for _, v := range balances {
		vv := v.(map[string]any)
		currency := NewCurrency(vv["asset"].(string), "").AdaptBccToBch()
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       num.ToFloat64(vv["free"]),
			ForzenAmount: num.ToFloat64(vv["locked"]),
		}
	}
	return &acc, nil
}
func (bn *Binance) LimitBuy(amount, price string, currencyPair CurrencyPair, opt ...LimitOrderOptionalParameter) (*q.Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "LIMIT", "BUY")
}
func (bn *Binance) LimitSell(amount, price string, currencyPair CurrencyPair, opt ...LimitOrderOptionalParameter) (*q.Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "LIMIT", "SELL")
}
func (bn *Binance) MarketBuy(amount, price string, currencyPair CurrencyPair) (*q.Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "MARKET", "BUY")
}
func (bn *Binance) MarketSell(amount, price string, currencyPair CurrencyPair) (*q.Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "MARKET", "SELL")
}
func (bn *Binance) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	path := bn.apiV3 + OrderUri
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	params.Set("orderId", orderId)
	bn.buildParamsSigned(&params)
	resp, err := web.HttpDeleteForm(bn.httpClient, path, params, bn.header())
	if err != nil {
		return false, bn.adaptError(err)
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return false, err
	}
	if orderIdCanceled := num.ToInt[int](respMap["orderId"]); orderIdCanceled <= 0 {
		return false, errors.New(slice.Bytes2String(resp))
	}
	return true, nil
}
func (bn *Binance) GetOneOrder(orderId string, currencyPair CurrencyPair) (*q.Order, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	if orderId != "" {
		params.Set("orderId", orderId)
	}
	params.Set("orderId", orderId)
	bn.buildParamsSigned(&params)
	path := bn.apiV3 + OrderUri + "?" + params.Encode()
	respMap, err := web.HttpGet2(bn.httpClient, path, bn.header())
	if err != nil {
		return nil, err
	}
	order := bn.adaptOrder(currencyPair, respMap)
	return &order, nil
}
func (bn *Binance) GetUnfinishedOrders(currencyPair CurrencyPair) ([]q.Order, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	bn.buildParamsSigned(&params)
	path := bn.apiV3 + UnfinishedOrdersInfo + params.Encode()
	respMap, err := web.HttpGet3(bn.httpClient, path, bn.header())
	if err != nil {
		return nil, err
	}
	orders := make([]q.Order, 0)
	for _, v := range respMap {
		ord := v.(map[string]any)
		orders = append(orders, bn.adaptOrder(currencyPair, ord))
	}
	return orders, nil
}
func (bn *Binance) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]Kline, error) {
	params := url.Values{}
	params.Set("symbol", currency.ToSymbol(""))
	params.Set("interval", internalKlinePeriodConverter[period])
	params.Set("limit", fmt.Sprintf("%d", size))
	util.MergeOptionalParameter(&params, optional...)
	klineUrl := bn.apiV3 + KlineUri + "?" + params.Encode()
	klines, err := web.HttpGet3(bn.httpClient, klineUrl, nil)
	if err != nil {
		return nil, err
	}
	var klineRecords []Kline
	for _, _record := range klines {
		r := Kline{Pair: currency}
		record := _record.([]any)
		r.Timestamp = int64(record[0].(float64)) / 1000 //to unix timestramp
		r.Open = num.ToFloat64(record[1])
		r.High = num.ToFloat64(record[2])
		r.Low = num.ToFloat64(record[3])
		r.Close = num.ToFloat64(record[4])
		r.Vol = num.ToFloat64(record[5])
		klineRecords = append(klineRecords, r)
	}
	return klineRecords, nil
}

// 非个人，整个交易所的交易记录
// 注意：since is fromId
func (bn *Binance) GetTrades(currencyPair CurrencyPair, since int64) ([]q.Trade, error) {
	param := url.Values{}
	param.Set("symbol", currencyPair.ToSymbol(""))
	param.Set("limit", "500")
	if since > 0 {
		param.Set("fromId", strconv.Itoa(int(since)))
	}
	apiUrl := bn.apiV3 + "historicalTrades?" + param.Encode()
	resp, err := web.HttpGet3(bn.httpClient, apiUrl, map[string]string{
		"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		return nil, err
	}
	var trades []q.Trade
	for _, v := range resp {
		m := v.(map[string]any)
		ty := SELL
		if m["isBuyerMaker"].(bool) {
			ty = BUY
		}
		trades = append(trades, q.Trade{
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
func (bn *Binance) GetOrderHistorys(currency CurrencyPair, optional ...OptionalParameter) ([]q.Order, error) {
	params := url.Values{}
	params.Set("symbol", currency.AdaptUsdToUsdt().ToSymbol(""))
	util.MergeOptionalParameter(&params, optional...)
	bn.buildParamsSigned(&params)
	path := bn.apiV3 + "allOrders?" + params.Encode()
	respMap, err := web.HttpGet3(bn.httpClient, path, bn.header())
	if err != nil {
		return nil, err
	}
	orders := make([]q.Order, 0)
	for _, v := range respMap {
		orderMap := v.(map[string]any)
		orders = append(orders, bn.adaptOrder(currency, orderMap))
	}
	return orders, nil
}
func (bn *Binance) toCurrencyPair(symbol string) CurrencyPair {
	if bn.ExchangeInfo == nil {
		var err error
		bn.ExchangeInfo, err = bn.GetExchangeInfo()
		if err != nil {
			return CurrencyPair{}
		}
	}
	for _, v := range bn.ExchangeInfo.Symbols {
		if v.Symbol == symbol {
			return NewCurrencyPair2(v.BaseAsset + "_" + v.QuoteAsset)
		}
	}
	return CurrencyPair{}
}
func (bn *Binance) GetExchangeInfo() (*ExchangeInfo, error) {
	resp, err := web.HttpGet5(bn.httpClient, bn.apiV3+"exchangeInfo", nil)
	if err != nil {
		return nil, err
	}
	info := &ExchangeInfo{}
	err = json.Unmarshal(resp, info)
	if err != nil {
		return nil, err
	}
	return info, nil
}
func (bn *Binance) GetTradeSymbol(currencyPair CurrencyPair) (*TradeSymbol, error) {
	if bn.ExchangeInfo == nil {
		var err error
		bn.ExchangeInfo, err = bn.GetExchangeInfo()
		if err != nil {
			return nil, err
		}
	}
	for _, v := range bn.ExchangeInfo.Symbols {
		if v.Symbol == currencyPair.ToSymbol("") {
			return &v, nil
		}
	}
	return nil, errors.New("symbol not found")
}
func (bn *Binance) adaptError(err error) error {
	errStr := err.Error()
	if strings.Contains(errStr, "Order does not exist") ||
		strings.Contains(errStr, "Unknown order sent") {
		return EX_ERR_NOT_FIND_ORDER.OriginErr(errStr)
	}
	if strings.Contains(errStr, "Too much request") {
		return EX_ERR_API_LIMIT.OriginErr(errStr)
	}
	if strings.Contains(errStr, "insufficient") {
		return EX_ERR_INSUFFICIENT_BALANCE.OriginErr(errStr)
	}
	return err
}
func (bn *Binance) adaptOrder(currencyPair CurrencyPair, orderMap map[string]any) q.Order {
	side := orderMap["side"].(string)
	orderSide := SELL
	if side == "BUY" {
		orderSide = BUY
	}
	quoteQty := num.ToFloat64(orderMap["cummulativeQuoteQty"])
	qty := num.ToFloat64(orderMap["executedQty"])
	avgPrice := 0.0
	if qty > 0 {
		avgPrice = num.FloatToFixed(quoteQty/qty, 8)
	}
	return q.Order{
		OrderID:      num.ToInt[int](orderMap["orderId"]),
		OrderID2:     fmt.Sprintf("%.0f", orderMap["orderId"]),
		Cid:          orderMap["clientOrderId"].(string),
		Currency:     currencyPair,
		Price:        num.ToFloat64(orderMap["price"]),
		Amount:       num.ToFloat64(orderMap["origQty"]),
		DealAmount:   num.ToFloat64(orderMap["executedQty"]),
		AvgPrice:     avgPrice,
		Side:         orderSide,
		Status:       adaptOrderStatus(orderMap["status"].(string)),
		OrderTime:    num.ToInt[int](orderMap["time"]),
		FinishedTime: num.ToInt[int64](orderMap["updateTime"]),
	}
}
