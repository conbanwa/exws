package gateio

import (
	"crypto/hmac"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/conbanwa/wstrader/q"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/conbanwa/logs"
)

const BaseUrl = "https://api.gateio.ws/api/v4"

func (g *Gate) PairArray() (sm map[string]q.D, ps map[q.D]q.P, err error) {
	sm = map[string]q.D{}
	ps = map[q.D]q.P{}
	type instrumentsResponse []struct {
		ID              string  `json:"id"`
		Base            string  `json:"base"`
		Quote           string  `json:"quote"`
		Fee             float64 `json:"fee,string"`
		MinBaseAmount   float64 `json:"min_base_amount,string"`
		MinQuoteAmount  float64 `json:"min_quote_amount,string"`
		AmountPrecision int     `json:"amount_precision"`
		Precision       int     `json:"precision"` //price
		TradeStatus     string  `json:"trade_status"`
		SellStart       int     `json:"sell_start"`
		BuyStart        int     `json:"buy_start"`
	}
	var response instrumentsResponse
	err = HttpGet4(g.client, BaseUrl+"/spot/currency_pairs", nil, &response)
	if err != nil {
		return
	}
	for _, v := range response {
		if v.TradeStatus == "tradable" {
			sm[v.ID] = Sym2duo(v.ID)
			ps[Sym2duo(v.ID)] = q.P{
				Base:         math.Pow10(-v.AmountPrecision),
				BaseLimit:    math.Pow10(-v.AmountPrecision),
				Quote:        math.Pow10(-v.Precision),
				Price:        math.Pow10(-v.Precision),
				MinBase:      v.MinBaseAmount,
				MinBaseLimit: v.MinBaseAmount,
				MinQuote:     v.MinQuoteAmount,
				TakerNett:    1 - v.Fee/100,
			}
		}
	}
	return
}
func Sym2duo(pair string) q.D {
	parts := strings.Split(pair, "_")
	var res q.D
	if len(parts) == 2 {
		res = q.D{Base: strings.ToUpper(parts[0]), Quote: strings.ToUpper(parts[1])}
	} else if pair == "timezone" {
		panic("FATAL: timezone" + pair)
	} else {
		panic("FATAL: DIV ERR!" + pair)
	}
	return res
}
func (g *Gate) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	orders = places
	for _, p := range places {
		go func(p q.Order) {
		}(p)
	}
	return
}
func (g *Gate) Test() bool {
	return true
}
func unify(local string) string {
	global := strings.ToUpper(local)
	return global
}
func (g *Gate) GetAttr() (a q.Attr) {
	a.MultiOrder = false
	return a
}

func (g *Gate) OneTicker(d q.D) (ticker q.Bbo, err error) {
	return
}
func (g *Gate) AllTickerV2(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	ret, err := HttpGet(g.client, BaseUrl+"/spot/tickers")
	if err != nil {
		logs.E(err)
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return mdt, errCode
	}
	for pairrev, v := range ret {
		pair := Sym2duo(pairrev)
		vm, ok := v.(map[string]any)
		if !ok {
			return mdt, errors.New("not found")
		}
		if vm["lowestAsk"] == nil {
			logs.I(vm)
			return mdt, fmt.Errorf("%+v", vm)
		}
		var ticker q.Bbo
		ticker.Pair = pairrev
		ticker.Bid, _ = strconv.ParseFloat(vm["highestBid"].(string), 64)
		ticker.Ask, _ = strconv.ParseFloat(vm["lowestAsk"].(string), 64)
		ticker.BidSize = 1e-12
		ticker.AskSize = 1e-12
		// ticker.High, _ = strconv.ParseFloat(vm["high24hr"].(string), 64)
		// ticker.Low, _ = strconv.ParseFloat(vm["low24hr"].(string), 64)
		// ticker.Last, _ = strconv.ParseFloat(vm["last"].(string), 64)
		// ticker.Vol, _ = strconv.ParseFloat(vm["quoteVolume"].(string), 64)
		if ticker.Valid() {
			mdt.Store(pair, ticker)
			// } else {
			// 	logs.E(ticker)
		}
	}
	if _, ok := mdt.Load(q.D{Base: "BTC", Quote: "USDT"}); !ok {
		logs.E(len(ret))
	}
	return
}
func (g *Gate) AllTicker(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	ret, err := HttpGet(g.client, marketBaseUrl+"/tickers")
	if err != nil {
		logs.E(err)
		errCode := HTTP_ERR_CODE
		errCode.OriginErrMsg = err.Error()
		return mdt, errCode
	}
	for pairrev, v := range ret {
		pairrev = strings.ToUpper(pairrev)
		pair, ok := SymPair[pairrev]
		if ok {
			vm, ok := v.(map[string]any)
			if !ok {
				return mdt, errors.New("not found")
			}
			if vm["lowestAsk"] == nil {
				logs.I(vm)
				return mdt, fmt.Errorf("%+v", vm)
			}
			var ticker q.Bbo
			ticker.Pair = pairrev
			ticker.Bid, _ = strconv.ParseFloat(vm["highestBid"].(string), 64)
			ticker.Ask, _ = strconv.ParseFloat(vm["lowestAsk"].(string), 64)
			ticker.BidSize = 1e-12
			ticker.AskSize = 1e-12
			if ticker.Valid() {
				mdt.Store(pair, ticker)
			}
			// }else{
			// 	logs.W(pairrev)
		}
	}
	return
}
func (g *Gate) httpDo(method string, url string, param string) []byte {
	key := []byte(g.secretkey)
	mac := hmac.New(sha512.New, key)
	mac.Write([]byte(param))
	var sign = fmt.Sprintf("%x", mac.Sum(nil))
	resp, err := NewRequest(g.client, "POST", url, "", map[string]string{
		"key":  g.accesskey,
		"sign": sign})
	if err != nil {
		logs.E(err)
	}
	return resp
}
