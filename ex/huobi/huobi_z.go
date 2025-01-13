package huobi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/conbanwa/exws/cons"
	"github.com/conbanwa/exws/q"
	"github.com/conbanwa/exws/stat/zelo"
	. "github.com/conbanwa/exws/util"
	. "github.com/conbanwa/exws/web"
	"github.com/conbanwa/num"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/conbanwa/logs"
)

var log = zelo.Writer.With().Str("ex", cons.HUOBI).Logger()

func tickerUrl(hbpro *HuoBiPro) string {
	return hbpro.baseUrl + "/market/tickers"
}
func duo2sym(d q.D) string {
	return strings.ToLower(d.Clip(""))
}
func (hb *HuoBiPro) PairArray() (map[string]q.D, map[q.D]q.P, error) {
	p := map[q.D]q.P{}
	sm := map[string]q.D{}
	Symbols, err := hb.GetCurrenciesPrecision()
	zelo.PanicOnErr(err).Send()
	// logs.D("huobi has ", len(Symbols), " symbols")
	for _, v := range Symbols {
		if v.Trading == "enabled" {
			d := q.D{Base: unify(v.BaseCurrency), Quote: unify(v.QuoteCurrency)}
			sm[v.Symbol] = d
			var symbol q.P
			symbol.Base = math.Pow10(-v.AmountPrecision)
			symbol.Quote = math.Pow10(-v.ValuePrecision)
			symbol.MinBase = v.MinAmount
			symbol.MinQuote = v.MinValue
			symbol.Price = math.Pow10(-v.PricePrecision)
			p[d] = symbol
		} else {
			logs.I("hb stop trading: ", v.Symbol)
		}
	}
	return sm, p, nil
}

func (hb *HuoBiPro) Fee() (f float64) {
	// f = 0.0005
	// return
	path := "/v2/reference/transact-fee-rate"
	params := &url.Values{}
	params.Set("symbols", "btcusdt")
	hb.buildPostForm("GET", path, params)
	ret, err := HttpGet(hb.httpClient, hb.baseUrl+path+"?"+params.Encode())
	if err != nil {
		logs.D(err)
		return
	}
	datai, ok := ret["data"]
	if !ok {
		logs.D(ret)
		return
	}
	data, ok := datai.([]any)
	if !ok {
		logs.D(ret["data"])
		return
	}
	fee, ok := data[0].(map[string]any)
	if !ok {
		logs.D(data)
		return
	}
	f, err = strconv.ParseFloat(fee["actualTakerRate"].(string), 64)
	if err != nil {
		logs.D(err)
	}
	return
}
func (hb *HuoBiPro) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	for i, p := range places {
		// go func(p Place) {
		sell := "buy-ioc"
		if p.Sell {
			sell = "sell-ioc"
		}
		n := p
		n.DealAmount = p.Amount
		// if p.Sell {
		// 	n.DealAmount /= p.Price
		// }
		orders[i] = n
		logs.W(p.Amount, p.Price, p.Symbol, sell)
		// orderId, err := hb.Place(FloatToString(p.Amount, p.AmountPrecision), FloatToString(p.Price, p.PricePrecision), p.Symbol, sell+"-ioc")
		// n.OrderID = num.ToInt[int](orderId)
		// n.OrderID2 = orderId
		// if err != nil {
		// 	return nil, err
		// }
		// }(p)
	}
	return
}
func (hb *HuoBiPro) Place(amount, price, symbol string, orderType string) (string, error) {
	path := "/v1/order/orders/place"
	params := url.Values{}
	params.Set("account-id", hb.accountId)
	params.Set("client-order-id", GenerateOrderClientId(32))
	params.Set("amount", amount)
	params.Set("symbol", symbol)
	params.Set("type", orderType)
	params.Set("price", price)
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
func (hb *HuoBiPro) Balances() (availables, frozens *sync.Map, err error) {
	availables, frozens = new(sync.Map), new(sync.Map)
	path := fmt.Sprintf("/v1/account/accounts/%s/balance", hb.accountId)
	params := &url.Values{}
	params.Set("accountId-id", hb.accountId)
	hb.buildPostForm("GET", path, params)
	urlStr := hb.baseUrl + path + "?" + params.Encode()
	//println(urlStr)
	respMap, err := HttpGet(hb.httpClient, urlStr)
	if err != nil {
		return
	}
	if respMap["status"].(string) != "ok" {
		err = fmt.Errorf(respMap["err-code"].(string))
		return
	}
	datamap := respMap["data"].(map[string]any)
	if datamap["state"].(string) != "working" {
		err = fmt.Errorf(respMap["state"].(string))
		return
	}
	list := datamap["list"].([]any)
	for _, v := range list {
		listmap := v.(map[string]any)
		currency := unify(listmap["currency"].(string))
		typeStr := listmap["type"].(string)
		balance := num.ToFloat64(listmap["balance"])
		switch typeStr {
		case "trade":
			availables.Store(currency, balance)
		case "frozen":
			frozens.Store(currency, balance)
		}
	}
	return
}
func (hb *HuoBiPro) Test() bool {
	return true
}
func unify(local string) string {
	global := strings.ToUpper(local)
	return global
}
func (hb *HuoBiPro) GetAttr() (a q.Attr) {
	a.MultiOrder = false
	return a
}
func (hb *HuoBiPro) TradeFee() (map[string]q.TradeFee, error) {
	return nil, nil
}
func (hb *HuoBiPro) WithdrawFee() (sf []q.NetworkWithdraw, err error) {

	return
}
func (hb *HuoBiPro) OneTicker(d q.D) (ticker q.Bbo, err error) {
	return
}
func (hb *HuoBiPro) AllTicker(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	ret, err := HttpGet(hb.httpClient, tickerUrl(hb))
	if err != nil {
		return
	}
	data, ok := ret["data"].([]any)
	if !ok {
		return mdt, errors.New("response format error")
	}
	for _, v := range data {
		vm := v.(map[string]any)
		if sym, ok := SymPair[vm["symbol"].(string)]; ok {
			var ticker q.Bbo
			ticker.Pair = vm["symbol"].(string)
			ticker.Bid = vm["bid"].(float64)
			ticker.BidSize = vm["bidSize"].(float64)
			ticker.Ask = vm["ask"].(float64)
			ticker.AskSize = vm["askSize"].(float64)
			// ticker.Amount = vm["amount"].(float64)
			if ticker.Valid() {
				mdt.Store(q.D{Base: unify(sym.Base), Quote: unify(sym.Quote)}, ticker)
			}
		}
	}
	return
}
