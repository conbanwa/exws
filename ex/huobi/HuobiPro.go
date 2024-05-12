package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/num"
	"github.com/conbanwa/slice"
	"math"
	"math/big"
	"net/http"
	"net/url"
	. "github.com/conbanwa/wstrader"
	. "github.com/conbanwa/wstrader/cons"
	. "github.com/conbanwa/wstrader/q"
	. "github.com/conbanwa/wstrader/util"
	. "github.com/conbanwa/wstrader/web"
	"sort"
	"strings"
	"time"

	"github.com/conbanwa/logs"
)

var HBPOINT = NewCurrency("HBPOINT", "")
var _INERNAL_KLINE_PERIOD_CONVERTER = map[KlinePeriod]string{
	KLINE_PERIOD_1MIN:   "1min",
	KLINE_PERIOD_5MIN:   "5min",
	KLINE_PERIOD_15MIN:  "15min",
	KLINE_PERIOD_30MIN:  "30min",
	KLINE_PERIOD_60MIN:  "60min",
	KLINE_PERIOD_1DAY:   "1day",
	KLINE_PERIOD_1WEEK:  "1week",
	KLINE_PERIOD_1MONTH: "1mon",
	KLINE_PERIOD_1YEAR:  "1year",
}

const (
	HB_POINT_ACCOUNT = "point"
	HB_SPOT_ACCOUNT  = "spot"
)

type AccountInfo struct {
	Id    string
	Type  string
	State string
}
type HuoBiPro struct {
	httpClient *http.Client
	baseUrl    string
	accountId  string
	accessKey  string
	secretKey  string
	Symbols    map[string]HuoBiProSymbol
	//ECDSAPrivateKey string
}
type HuoBiProSymbol struct {
	BaseCurrency    string
	QuoteCurrency   string
	PricePrecision  int
	AmountPrecision int
	ValuePrecision  int
	MinAmount       float64
	MinValue        float64
	SymbolPartition string
	Symbol          string
	Trading         string
}

func NewHuobiWithConfig(config *APIConfig) *HuoBiPro {
	hbpro := new(HuoBiPro)
	if config.Endpoint == "" {
		hbpro.baseUrl = "https://api.huobi.pro"
	} else {
		hbpro.baseUrl = config.Endpoint
	}
	hbpro.httpClient = config.HttpClient
	hbpro.accessKey = config.ApiKey
	hbpro.secretKey = config.ApiSecretKey
	if config.ApiKey != "" && config.ApiSecretKey != "" {
		accinfo, err := hbpro.GetAccountInfo(HB_SPOT_ACCOUNT)
		if err != nil {
			hbpro.accountId = ""
			//panic(err)
		} else {
			hbpro.accountId = accinfo.Id
			//log.Println("account state :", accinfo.State)
			logs.I("accountId=", accinfo.Id, ",state=", accinfo.State, ",type=", accinfo.Type)
		}
	}
	hbpro.Symbols = make(map[string]HuoBiProSymbol, 100)
	_, err := hbpro.GetCurrenciesPrecision()
	if err != nil {
		panic("GetCurrenciesPrecision Error=" + err.Error())
	}
	return hbpro
}
func NewHuoBiPro(client *http.Client, apikey, secretkey, accountId string) *HuoBiPro {
	hbpro := new(HuoBiPro)
	hbpro.baseUrl = "https://api.huobi.pro"
	hbpro.httpClient = client
	hbpro.accessKey = apikey
	hbpro.secretKey = secretkey
	hbpro.accountId = accountId
	return hbpro
}

/**
 *现货交易
 */
func NewHuoBiProSpot(client *http.Client, apikey, secretkey string) *HuoBiPro {
	hb := NewHuoBiPro(client, apikey, secretkey, "")
	accinfo, err := hb.GetAccountInfo(HB_SPOT_ACCOUNT)
	if err != nil {
		hb.accountId = ""
		panic(err)
	} else {
		hb.accountId = accinfo.Id
		logs.I("account state :", accinfo.State)
	}
	hb.Symbols = make(map[string]HuoBiProSymbol, 100)
	_, err = hb.GetCurrenciesPrecision()
	if err != nil {
		panic("GetCurrenciesPrecision Error=" + err.Error())
	}
	return hb
}

/**
 * 点卡账户
 */
func NewHuoBiProPoint(client *http.Client, apikey, secretkey string) *HuoBiPro {
	hb := NewHuoBiPro(client, apikey, secretkey, "")
	accinfo, err := hb.GetAccountInfo(HB_POINT_ACCOUNT)
	if err != nil {
		panic(err)
	}
	hb.accountId = accinfo.Id
	logs.I("account state :" + accinfo.State)
	return hb
}
func (hb *HuoBiPro) GetAccountInfo(acc string) (AccountInfo, error) {
	path := "/v1/account/accounts"
	params := &url.Values{}
	hb.buildPostForm("GET", path, params)
	//log.Println(hb.baseUrl + path + "?" + params.Encode())
	respMap, err := HttpGet(hb.httpClient, hb.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return AccountInfo{}, err
	}
	if respMap["status"].(string) != "ok" {
		return AccountInfo{}, errors.New(respMap["err-code"].(string))
	}
	var info AccountInfo
	data := respMap["data"].([]any)
	for _, v := range data {
		iddata := v.(map[string]any)
		if iddata["type"].(string) == acc {
			info.Id = fmt.Sprintf("%.0f", iddata["id"])
			info.Type = acc
			info.State = iddata["state"].(string)
			break
		}
	}
	//log.Println(respMap)
	return info, nil
}
func (hb *HuoBiPro) GetAccount() (*Account, error) {
	path := fmt.Sprintf("/v1/account/accounts/%s/balance", hb.accountId)
	params := &url.Values{}
	params.Set("accountId-id", hb.accountId)
	hb.buildPostForm("GET", path, params)
	urlStr := hb.baseUrl + path + "?" + params.Encode()
	logs.D(hb.accessKey)
	respMap, err := HttpGet(hb.httpClient, urlStr)
	if err != nil {
		return nil, err
	}
	//log.Println(respMap)
	if respMap["status"].(string) != "ok" {
		return nil, errors.New(respMap["err-code"].(string))
	}
	datamap := respMap["data"].(map[string]any)
	if datamap["state"].(string) != "working" {
		return nil, errors.New(datamap["state"].(string))
	}
	list := datamap["list"].([]any)
	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount, 6)
	acc.Exchange = hb.String()
	subAccMap := make(map[Currency]*SubAccount)
	for _, v := range list {
		balancemap := v.(map[string]any)
		currencySymbol := balancemap["currency"].(string)
		currency := NewCurrency(currencySymbol, "")
		typeStr := balancemap["type"].(string)
		balance := num.ToFloat64(balancemap["balance"])
		if subAccMap[currency] == nil {
			subAccMap[currency] = new(SubAccount)
		}
		subAccMap[currency].Currency = currency
		switch typeStr {
		case "trade":
			subAccMap[currency].Amount = balance
		case "frozen":
			subAccMap[currency].ForzenAmount = balance
		}
	}
	for k, v := range subAccMap {
		acc.SubAccounts[k] = *v
	}
	return acc, nil
}
func (hb *HuoBiPro) placeOrder(amount, price string, pair CurrencyPair, orderType string) (string, error) {
	symbol := hb.Symbols[pair.ToLower().ToSymbol("")]
	path := "/v1/order/orders/place"
	params := url.Values{}
	params.Set("account-id", hb.accountId)
	params.Set("client-order-id", GenerateOrderClientId(32))
	params.Set("amount", num.FloatToString(num.ToFloat64(amount), math.Pow10(-symbol.AmountPrecision)))
	params.Set("symbol", pair.AdaptUsdToUsdt().ToLower().ToSymbol(""))
	params.Set("type", orderType)
	switch orderType {
	case "buy-limit", "sell-limit":
		params.Set("price", num.FloatToString(num.ToFloat64(price), math.Pow10(-symbol.PricePrecision)))
	}
	hb.buildPostForm("POST", path, &params)
	resp, err := HttpPostForm3(hb.httpClient, hb.baseUrl+path+"?"+params.Encode(), hb.toJson(params),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	if err != nil {
		return "", err
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return "", err
	}
	if respMap["status"].(string) != "ok" {
		return "", errors.New(respMap["err-code"].(string))
	}
	return respMap["data"].(string), nil
}
func (hb *HuoBiPro) LimitBuy(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	orderTy := "buy-limit"
	if len(opt) > 0 {
		switch opt[0] {
		case PostOnly:
			orderTy = "buy-limit-maker"
		case Ioc:
			orderTy = "buy-ioc"
		case Fok:
			orderTy = "buy-limit-fok"
		default:
			logs.E("limit order optional parameter error ,opt= ", opt[0])
		}
	}
	orderId, err := hb.placeOrder(amount, price, currency, orderTy)
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  num.ToInt[int](orderId),
		OrderID2: orderId,
		Amount:   num.ToFloat64(amount),
		Price:    num.ToFloat64(price),
		Side:     BUY}, nil
}
func (hb *HuoBiPro) LimitSell(amount, price string, currency CurrencyPair, opt ...LimitOrderOptionalParameter) (*Order, error) {
	orderTy := "sell-limit"
	if len(opt) > 0 {
		switch opt[0] {
		case PostOnly:
			orderTy = "sell-limit-maker"
		case Ioc:
			orderTy = "sell-ioc"
		case Fok:
			orderTy = "sell-limit-fok"
		default:
			logs.E("limit order optional parameter error ,opt= ", opt[0])
		}
	}
	orderId, err := hb.placeOrder(amount, price, currency, orderTy)
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  num.ToInt[int](orderId),
		OrderID2: orderId,
		Amount:   num.ToFloat64(amount),
		Price:    num.ToFloat64(price),
		Side:     SELL}, nil
}
func (hb *HuoBiPro) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hb.placeOrder(amount, price, currency, "buy-market")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  num.ToInt[int](orderId),
		OrderID2: orderId,
		Amount:   num.ToFloat64(amount),
		Price:    num.ToFloat64(price),
		Side:     BUY_MARKET}, nil
}
func (hb *HuoBiPro) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	orderId, err := hb.placeOrder(amount, price, currency, "sell-market")
	if err != nil {
		return nil, err
	}
	return &Order{
		Currency: currency,
		OrderID:  num.ToInt[int](orderId),
		OrderID2: orderId,
		Amount:   num.ToFloat64(amount),
		Price:    num.ToFloat64(price),
		Side:     SELL_MARKET}, nil
}
func (hb *HuoBiPro) parseOrder(ordmap map[string]any) Order {
	ord := Order{
		Cid:        fmt.Sprint(ordmap["client-order-id"]),
		OrderID:    num.ToInt[int](ordmap["id"]),
		OrderID2:   fmt.Sprint(num.ToInt[int](ordmap["id"])),
		Amount:     num.ToFloat64(ordmap["amount"]),
		Price:      num.ToFloat64(ordmap["price"]),
		DealAmount: num.ToFloat64(ordmap["field-amount"]),
		Fee:        num.ToFloat64(ordmap["field-fees"]),
		OrderTime:  num.ToInt[int](ordmap["created-at"]),
	}
	state := ordmap["state"].(string)
	switch state {
	case "submitted", "pre-submitted":
		ord.Status = ORDER_UNFINISH
	case "filled":
		ord.Status = ORDER_FINISH
	case "partial-filled":
		ord.Status = ORDER_PART_FINISH
	case "canceled", "partial-canceled":
		ord.Status = ORDER_CANCEL
	default:
		ord.Status = ORDER_UNFINISH
	}
	if ord.DealAmount > 0.0 {
		ord.AvgPrice = num.ToFloat64(ordmap["field-cash-amount"]) / ord.DealAmount
	}
	typeS := ordmap["type"].(string)
	switch typeS {
	case "buy-limit":
		ord.Side = BUY
	case "buy-market":
		ord.Side = BUY_MARKET
	case "sell-limit":
		ord.Side = SELL
	case "sell-market":
		ord.Side = SELL_MARKET
	}
	return ord
}
func (hb *HuoBiPro) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	path := "/v1/order/orders/" + orderId
	params := url.Values{}
	hb.buildPostForm("GET", path, &params)
	respMap, err := HttpGet(hb.httpClient, hb.baseUrl+path+"?"+params.Encode())
	if err != nil {
		return nil, err
	}
	if respMap["status"].(string) != "ok" {
		return nil, errors.New(respMap["err-code"].(string))
	}
	datamap := respMap["data"].(map[string]any)
	order := hb.parseOrder(datamap)
	order.Currency = currency
	return &order, nil
}
func (hb *HuoBiPro) GetUnfinishedOrders(currency CurrencyPair) ([]Order, error) {
	return hb.getOrders(currency, OptionalParameter{}.
		Optional("states", "pre-submitted,submitted,partial-filled").
		Optional("size", "100"))
}
func (hb *HuoBiPro) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	path := fmt.Sprintf("/v1/order/orders/%s/submitcancel", orderId)
	params := url.Values{}
	hb.buildPostForm("POST", path, &params)
	resp, err := HttpPostForm3(hb.httpClient, hb.baseUrl+path+"?"+params.Encode(), hb.toJson(params),
		map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"})
	if err != nil {
		return false, err
	}
	var respMap map[string]any
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		return false, err
	}
	if respMap["status"].(string) != "ok" {
		return false, errors.New(slice.Bytes2String(resp))
	}
	return true, nil
}
func (hb *HuoBiPro) GetOrderHistorys(currency CurrencyPair, optional ...OptionalParameter) ([]Order, error) {
	var optionals []OptionalParameter
	optionals = append(optionals, OptionalParameter{}.
		Optional("states", "canceled,partial-canceled,filled").
		Optional("size", "100").
		Optional("direct", "next"))
	optionals = append(optionals, optional...)
	return hb.getOrders(currency, optionals...)
}

type queryOrdersParams struct {
	types,
	startDate,
	endDate,
	states,
	from,
	direct string
	size int
	pair CurrencyPair
}

func (hb *HuoBiPro) getOrders(pair CurrencyPair, optional ...OptionalParameter) ([]Order, error) {
	path := "/v1/order/orders"
	params := url.Values{}
	params.Set("symbol", strings.ToLower(pair.AdaptUsdToUsdt().ToSymbol("")))
	MergeOptionalParameter(&params, optional...)
	logs.I(params)
	hb.buildPostForm("GET", path, &params)
	respMap, err := HttpGet(hb.httpClient, fmt.Sprintf("%s%s?%s", hb.baseUrl, path, params.Encode()))
	if err != nil {
		return nil, err
	}
	if respMap["status"].(string) != "ok" {
		return nil, errors.New(respMap["err-code"].(string))
	}
	datamap := respMap["data"].([]any)
	var orders []Order
	for _, v := range datamap {
		ordmap := v.(map[string]any)
		ord := hb.parseOrder(ordmap)
		ord.Currency = pair
		orders = append(orders, ord)
	}
	return orders, nil
}
func (hb *HuoBiPro) GetTicker(currencyPair CurrencyPair) (*Ticker, error) {
	pair := currencyPair.AdaptUsdToUsdt()
	uri := hb.baseUrl + "/market/detail/merged?symbol=" + strings.ToLower(pair.ToSymbol(""))
	resp, err := HttpGet(hb.httpClient, uri)
	if err != nil {
		return nil, err
	}
	if resp["status"].(string) == "error" {
		return nil, errors.New(resp["err-msg"].(string))
	}
	tick, ok := resp["tick"].(map[string]any)
	if !ok {
		return nil, errors.New("tick assert error")
	}
	ticker := new(Ticker)
	ticker.Pair = currencyPair
	ticker.Vol = num.ToFloat64(tick["amount"])
	ticker.Low = num.ToFloat64(tick["low"])
	ticker.High = num.ToFloat64(tick["high"])
	bid, isOk := tick["bid"].([]any)
	if isOk != true {
		return nil, errors.New("no bid")
	}
	ask, isOk := tick["ask"].([]any)
	if isOk != true {
		return nil, errors.New("no ask")
	}
	ticker.Buy = num.ToFloat64(bid[0])
	ticker.Sell = num.ToFloat64(ask[0])
	ticker.Last = num.ToFloat64(tick["close"])
	ticker.Date = num.ToInt[uint64](resp["ts"])
	return ticker, nil
}
func (hb *HuoBiPro) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	uri := hb.baseUrl + "/market/depth?symbol=%s&type=step0&depth=%d"
	n := 5
	pair := currency.AdaptUsdToUsdt()
	if size <= 5 {
		n = 5
	} else if size <= 10 {
		n = 10
	} else if size <= 20 {
		n = 20
	} else {
		uri = hb.baseUrl + "/market/depth?symbol=%s&type=step0&d=%d"
	}
	respMap, err := HttpGet(hb.httpClient, fmt.Sprintf(uri, strings.ToLower(pair.ToSymbol("")), n))
	if err != nil {
		return nil, err
	}
	if "ok" != respMap["status"].(string) {
		return nil, errors.New(respMap["err-msg"].(string))
	}
	tick, _ := respMap["tick"].(map[string]any)
	dep := hb.parseDepthData(tick, size)
	dep.Pair = currency
	mills := num.ToInt[uint64](tick["ts"])
	dep.UTime = time.Unix(int64(mills/1000), int64(mills%1000)*int64(time.Millisecond))
	return dep, nil
}

// 倒序
func (hb *HuoBiPro) GetKlineRecords(currency CurrencyPair, period KlinePeriod, size int, optional ...OptionalParameter) ([]Kline, error) {
	uri := hb.baseUrl + "/market/history/kline?period=%s&size=%d&symbol=%s"
	symbol := strings.ToLower(currency.AdaptUsdToUsdt().ToSymbol(""))
	periodS, isOk := _INERNAL_KLINE_PERIOD_CONVERTER[period]
	if isOk != true {
		periodS = "1min"
	}
	ret, err := HttpGet(hb.httpClient, fmt.Sprintf(uri, periodS, size, symbol))
	if err != nil {
		return nil, err
	}
	data, ok := ret["data"].([]any)
	if !ok {
		return nil, errors.New("response format error")
	}
	var klines []Kline
	for _, e := range data {
		item := e.(map[string]any)
		klines = append(klines, Kline{
			Pair:      currency,
			Open:      num.ToFloat64(item["open"]),
			Close:     num.ToFloat64(item["close"]),
			High:      num.ToFloat64(item["high"]),
			Low:       num.ToFloat64(item["low"]),
			Vol:       num.ToFloat64(item["amount"]),
			Timestamp: num.ToInt[int64](item["id"])})
	}
	return klines, nil
}
func (hb *HuoBiPro) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	var (
		trades []Trade
		ret    struct {
			Status string
			ErrMsg string `json:"err-msg"`
			Data   []struct {
				Ts   int64
				Data []struct {
					Id        big.Int
					Amount    float64
					Price     float64
					Direction string
					Ts        int64
				}
			}
		}
	)
	uri := hb.baseUrl + "/market/history/trade?size=2000&symbol=" + currencyPair.AdaptUsdToUsdt().ToLower().ToSymbol("")
	err := HttpGet4(hb.httpClient, uri, map[string]string{}, &ret)
	if err != nil {
		return nil, err
	}
	if ret.Status != "ok" {
		return nil, errors.New(ret.ErrMsg)
	}
	for _, d := range ret.Data {
		for _, t := range d.Data {
			//fix huobi   Weird rules of tid
			//火币交易ID规定固定23位, 导致超出int64范围，每个交易对有不同的固定填充前缀
			//实际交易ID远远没有到23位数字。
			tid := num.ToInt[int64](strings.TrimPrefix(t.Id.String()[4:], "0"))
			if tid == 0 {
				tid = num.ToInt[int64](strings.TrimPrefix(t.Id.String()[5:], "0"))
			}
			///
			trades = append(trades, Trade{
				Tid:    num.ToInt[int64](tid),
				Pair:   currencyPair,
				Amount: t.Amount,
				Price:  t.Price,
				Type:   AdaptTradeSide(t.Direction),
				Date:   t.Ts})
		}
	}
	return trades, nil
}

type ecdsaSignature struct {
	R, S *big.Int
}

func (hb *HuoBiPro) buildPostForm(reqMethod, path string, postForm *url.Values) error {
	postForm.Set("AccessKeyId", hb.accessKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))
	domain := strings.Replace(hb.baseUrl, "https://", "", len(hb.baseUrl))
	payload := fmt.Sprintf("%s\n%s\n%s\n%s", reqMethod, domain, path, postForm.Encode())
	sign, _ := GetParamHmacSHA256Base64Sign(hb.secretKey, payload)
	postForm.Set("Signature", sign)
	/**
	p, _ := pem.Decode([]byte(hb.ECDSAPrivateKey))
	pri, _ := secp256k1_go.PrivKeyFromBytes(secp256k1_go.S256(), p.Bytes)
	signer, _ := pri.Sign([]byte(sign))
	signAsn, _ := asn1.Marshal(signer)
	priSign := base64.StdEncoding.EncodeToString(signAsn)
	postForm.Set("PrivateSignature", priSign)
	*/
	return nil
}
func (hb *HuoBiPro) toJson(params url.Values) string {
	parammap := make(map[string]string)
	for k, v := range params {
		parammap[k] = v[0]
	}
	jsonData, _ := json.Marshal(parammap)
	return slice.Bytes2String(jsonData)
}
func (hb *HuoBiPro) parseDepthData(tick map[string]any, size int) *Depth {
	bids, _ := tick["bids"].([]any)
	asks, _ := tick["asks"].([]any)
	depth := new(Depth)
	n := 0
	for _, r := range asks {
		var dr DepthRecord
		rr := r.([]any)
		dr.Price = num.ToFloat64(rr[0])
		dr.Amount = num.ToFloat64(rr[1])
		depth.AskList = append(depth.AskList, dr)
		n++
		if n == size {
			break
		}
	}
	n = 0
	for _, r := range bids {
		var dr DepthRecord
		rr := r.([]any)
		dr.Price = num.ToFloat64(rr[0])
		dr.Amount = num.ToFloat64(rr[1])
		depth.BidList = append(depth.BidList, dr)
		n++
		if n == size {
			break
		}
	}
	sort.Sort(sort.Reverse(depth.AskList))
	return depth
}
func (hb *HuoBiPro) String() string {
	return HUOBI
}
func (hb *HuoBiPro) GetCurrenciesList() ([]string, error) {
	uri := hb.baseUrl + "/v1/common/currencys"
	ret, err := HttpGet(hb.httpClient, uri)
	if err != nil {
		return nil, err
	}
	data, ok := ret["data"].([]any)
	if !ok {
		return nil, errors.New("response format error")
	}
	logs.D(data)
	return nil, nil
}
func (hb *HuoBiPro) GetCurrenciesPrecision() ([]HuoBiProSymbol, error) {
	uri := hb.baseUrl + "/v1/common/symbols"
	ret, err := HttpGet(hb.httpClient, uri)
	if err != nil {
		return nil, err
	}
	data, ok := ret["data"].([]any)
	if !ok {
		return nil, errors.New("response format error")
	}
	var Symbols []HuoBiProSymbol
	for _, v := range data {
		_sym := v.(map[string]any)
		var sym HuoBiProSymbol
		sym.BaseCurrency = _sym["base-currency"].(string)
		sym.QuoteCurrency = _sym["quote-currency"].(string)
		sym.PricePrecision = int(_sym["price-precision"].(float64))
		sym.AmountPrecision = int(_sym["amount-precision"].(float64))
		sym.ValuePrecision = int(_sym["value-precision"].(float64))
		sym.MinAmount = _sym["limit-order-min-order-amt"].(float64)
		sym.MinValue = _sym["min-order-value"].(float64)
		sym.SymbolPartition = _sym["symbol-partition"].(string)
		sym.Symbol = _sym["symbol"].(string)
		sym.Trading = _sym["api-trading"].(string)
		Symbols = append(Symbols, sym)
		hb.Symbols[sym.Symbol] = sym
	}
	//logs.D(Symbols)
	return Symbols, nil
}
