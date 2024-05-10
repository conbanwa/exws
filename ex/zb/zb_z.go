package zb

import (
	"encoding/json"
	"fmt"
	"github.com/conbanwa/num"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/conbanwa/logs"
)

const turl = "https://api.zb.com/data/v1/" //zb.today
var SymPair = make(map[string]q.D)

func (zb *Zb) PairArray() (map[string]q.D, map[q.D]q.P, error) {
	var sm = map[string]q.D{}
	p := map[q.D]q.P{}
	resp, err := HttpZBP(zb.httpClient, turl+"markets")
	if err != nil {
		return nil, nil, err
	}
	// logs.E(resp)
	for sym, tickermap := range resp {
		var pre q.P
		pre.Base = math.Pow10(int(-tickermap["amountScale"]))
		pre.Quote = math.Pow10(-int(tickermap["priceScale"]))
		pre.Price = math.Pow10(-int(tickermap["priceScale"]))
		sm[PairNoLine(sym)] = Sym2duo(sym)
		p[Sym2duo(sym)] = pre
		SymPair[PairNoLine(sym)] = Sym2duo(sym)
		// logs.F(sym, Sym2duo(sym),PairNoLine(sym))
	}
	if len(sm) == 0 {
		panic("FATAL: no symap")
	}
	return sm, p, nil
}
func duo2sym(d q.D) string {
	return strings.ToLower(d.Clip("_"))
}
func Sym2duo(pair string) q.D {
	parts := strings.Split(pair, "_")
	var res q.D
	if len(parts) == 2 {
		res = q.D{Base: unify(parts[0]), Quote: unify(parts[1])}
	} else {
		logs.F("FATAL: DIV ERR!", pair)
	}
	if res.Base == res.Quote {
		logs.F("FATAL: SAME CURRENCY ERR!", pair)
	}
	return res
}
func PairNoLine(pair string) (res string) {
	parts := strings.Split(pair, "_")
	if len(parts) == 2 {
		res = parts[0] + parts[1]
	} else {
		logs.F("FATAL: DIV ERR!", pair)
	}
	return
}
func (zb *Zb) Fee() float64 {
	return 0.002
}
func (zb *Zb) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	orders = places
	for _, p := range places {
		go func(p q.Order) {
		}(p)
	}
	return
}
func (zb *Zb) Balances() (availables, frozens *sync.Map, err error) {
	availables, frozens = new(sync.Map), new(sync.Map)
	params := url.Values{}
	params.Set("method", "getAccountInfo")
	zb.buildPostForm(&params)
	//logs.E(params.Encode())
	resp, err := HttpPostForm(zb.httpClient, TradeUrl+GetAccountApi, params)
	if err != nil {
		logs.E(err)
		return
	}
	var respMap map[string]any
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		logs.E("json unmarshal error")
		logs.E(err)
		return
	}
	if respMap["code"] != nil && respMap["code"].(float64) != 1000 {
		// logs.E(string(resp))
		return
	}
	acc := new(Account)
	acc.Exchange = zb.String()
	acc.SubAccounts = make(map[cons.Currency]SubAccount)
	resultmap := respMap["result"].(map[string]any)
	coins := resultmap["coins"].([]any)
	acc.NetAsset = num.ToFloat64(resultmap["netAssets"])
	acc.Asset = num.ToFloat64(resultmap["totalAssets"])
	for _, v := range coins {
		vv := v.(map[string]any)
		currency := unify(vv["key"].(string))
		availables.Store(currency, num.ToFloat64(vv["available"]))
		frozens.Store(currency, num.ToFloat64(vv["freez"]))
	}
	return
}
func (zb *Zb) Test() bool {
	return true
}
func unify(local string) string {
	global := strings.ToUpper(local)
	return global
}
func (zb *Zb) GetAttr() (a q.Attr) {
	a.MultiOrder = true
	return a
}

func (zb *Zb) TradeFee() (map[string]q.TradeFee, error) {
	return nil, nil
}
func (zb *Zb) WithdrawFee() (sf []q.NetworkWithdraw, err error) {

	return
}
func (zb *Zb) OneTicker(d q.D) (ticker q.Bbo, err error) {
	resp, err := HttpGet(zb.httpClient, MarketUrl+fmt.Sprintf(DepthApi, duo2sym(d), 1))
	if err != nil {
		return
	}
	asks, isok1 := resp["asks"].([]any)
	bids, isok2 := resp["bids"].([]any)
	if !isok2 || !isok1 {
		err = fmt.Errorf("no depth data")
		return
	}
	for _, e := range bids {
		ee := e.([]any)
		ticker.Bid = ee[0].(float64)
		ticker.BidSize = ee[1].(float64)
	}
	for _, e := range asks {
		ee := e.([]any)
		ticker.Ask = ee[0].(float64)
		ticker.AskSize = ee[1].(float64)
	}
	return
}
func (zb *Zb) AllTicker(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	if len(SymPair) == 0 {
		panic("FATAL: no symap")
	}
	resp, err := HttpGet(zb.httpClient, turl+"allTicker")
	if err != nil {
		return
	}
	//logs.E(resp)
	for sym, v := range resp {
		tickermap := v.(map[string]any)
		var ticker q.Bbo
		ticker.Bid, _ = strconv.ParseFloat(tickermap["buy"].(string), 64)
		ticker.Ask, _ = strconv.ParseFloat(tickermap["sell"].(string), 64)
		ticker.BidSize = 1e-12
		ticker.AskSize = 1e-12
		ticker.Pair = sym
		if !ticker.Valid() {
			// logs.E("ZB ", sym, ticker, resp)
		} else {
			if SymPair[sym].Valid() {
				mdt.Store(SymPair[sym], ticker)
			}
		}
	}
	return
}
