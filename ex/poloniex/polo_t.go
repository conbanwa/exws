package poloniex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	. "github.com/conbanwa/wstrader"
	"github.com/conbanwa/wstrader/cons"
	"github.com/conbanwa/wstrader/q"
	. "github.com/conbanwa/wstrader/web"
	"strconv"
	"strings"
	"sync"

	"github.com/conbanwa/logs"
)

func (poloniex *Poloniex) PairArray() (map[string]q.D, map[q.D]q.P, error) {
	tickerUrl := PUBLIC_URL + TICKER_API
	ret, err := HttpGet(poloniex.client, tickerUrl)
	if err != nil {
		return nil, nil, err
	}
	var sm = map[string]q.D{}
	var ps = map[q.D]q.P{}
	for pairrev := range ret {
		sm[pairrev] = Sym2duo(pairrev)
		ps[Sym2duo(pairrev)] = q.DefaultPrecision()
	}
	return sm, ps, nil
}
func Sym2duo(pair string) q.D {
	parts := strings.Split(pair, "_")
	var res q.D
	if len(parts) == 2 {
		res = q.D{Base: parts[1], Quote: parts[0]}
	} else if pair == "timezone" {
		logs.F("FATAL: timezone", pair)
	} else {
		logs.F("FATAL: DIV ERR!", pair)
	}
	return res
}
func (poloniex *Poloniex) Fee() float64 {
	// url := PUBLIC_URL + "TICKER_API"
	params := url.Values{}
	params.Set("command", "returnFeeInfo")
	resp, err := poloniex.Post(params)
	if err != nil {
		return 0.002
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil {
		logs.D(respMap)
		return 0.002
	}
	f, err := strconv.ParseFloat(respMap["takerFee"].(string), 64)
	if err != nil {
		logs.D(respMap)
		return 0.002
	}
	return f
}
func (poloniex *Poloniex) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	for _, p := range places {
		go func(p q.Order) {
			// 	price := p.Price
			// 	amount := p.Amount
			// 	command := "sell"
			// 	currency := Clip(p.D), "_".Reverse()
			// 	if p.Sell {
			// 		command = "buy"
			// 	}
			// 	postData := url.Values{}
			// 	postData.Set("command", command)
			// 	postData.Set("currencyPair", currency)
			// 	postData.Set("rate", strconv.FormatFloat(price, 'f', 10, 64))
			// 	postData.Set("amount", strconv.FormatFloat(amount, 'f', 10, 64))
			// 	sign, _ := poloniex.buildPostForm(&postData)
			// 	headers := map[string]string{
			// 		"Key":  poloniex.accessKey,
			// 		"Sign": sign}
			// 	resp, err := HttpPostForm2(poloniex.client, TRADE_API, postData, headers)
			// 	if err != nil {
			// 		logs.E(err)
			// 		return nil, err
			// 	}
			// respMap := make(map[string]any)
			// err = json.Unmarshal(resp, &respMap)
			// if err != nil || respMap["error"] != nil {
			// 	logs.I(err, string(resp))
			// 	return nil, err
			// }
			// orderNumber := respMap["orderNumber"].(string)
			// order := new(Order)
			// order.OrderID, _ = strconv.Atoi(orderNumber)
			// order.OrderID2 = orderNumber
			// order.Amount, _ = strconv.ParseFloat(amount, 64)
			// order.Price, _ = strconv.ParseFloat(price, 64)
			// order.Status = ORDER_UNFINISH
			// order.Currency = currency
			// switch command {
			// case "sell":
			// 	order.Side = SELL
			// case "buy":
			// 	order.Side = BUY
			// }
		}(p)
	}
	return
}
func (poloniex *Poloniex) Balances() (availables, frozens *sync.Map, err error) {
	availables, frozens = new(sync.Map), new(sync.Map)
	postData := url.Values{}
	postData.Add("command", "returnCompleteBalances")
	sign, err := poloniex.buildPostForm(&postData)
	if err != nil {
		logs.E(err)
		return
	}
	headers := map[string]string{
		"Key":  poloniex.accessKey,
		"Sign": sign}
	resp, err := HttpPostForm2(poloniex.client, TRADE_API, postData, headers)
	if err != nil {
		logs.E(err)
		return
	}
	respMap := make(map[string]any)
	err = json.Unmarshal(resp, &respMap)
	if err != nil || respMap["error"] != nil {
		logs.I(err, respMap["error"])
		return
	}
	acc := new(Account)
	acc.Exchange = EXCHANGE_NAME
	acc.SubAccounts = make(map[cons.Currency]SubAccount)
	for k, v := range respMap {
		vv := v.(map[string]any)
		ava, _ := strconv.ParseFloat(vv["available"].(string), 64)
		availables.Store(k, ava)
		fro, _ := strconv.ParseFloat(vv["onOrders"].(string), 64)
		frozens.Store(k, fro)
	}
	return
}
func (poloniex *Poloniex) Test() bool {
	return true
}
func unify(local string) string {
	global := strings.ToUpper(local)
	return global
}
func (poloniex *Poloniex) GetAttr() (a q.Attr) {
	a.MultiOrder = false
	return a
}

func (poloniex *Poloniex) TradeFee() (map[string]q.TradeFee, error) {
	return nil, nil
}
func (poloniex *Poloniex) WithdrawFee() (sf []q.NetworkWithdraw, err error) {

	return
}
func (poloniex *Poloniex) OneTicker(d q.D) (ticker q.Bbo, err error) {
	return
}
func (poloniex *Poloniex) AllTicker(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	ret, err := HttpGet(poloniex.client, PUBLIC_URL+
		fmt.Sprintf(ORDER_BOOK_API, "all", 1))
	if err != nil {
		return
	}
	for pairrev, v := range ret {
		pair := Sym2duo(pairrev)
		tickermap, ok := v.(map[string]any)
		if !ok {
			return mdt, errors.New("not found")
		}
		isFrozen, o1 := tickermap["isFrozen"].(string)
		// postOnly, o2 := tickermap["postOnly"].(string)
		if isFrozen == "0" && o1 {
			if tickermap["asks"] == nil {
				logs.I(tickermap)
				return mdt, fmt.Errorf("%+v", tickermap)
			}
			if _, isOK := tickermap["asks"].([]any); !isOK {
				logs.I(tickermap)
				return mdt, fmt.Errorf("%+v", tickermap)
			}
			var ticker q.Bbo
			for _, v := range tickermap["asks"].([]any) {
				for i, vv := range v.([]any) {
					switch i {
					case 0:
						ticker.Ask, _ = strconv.ParseFloat(vv.(string), 64)
					case 1:
						ticker.AskSize = vv.(float64)
					}
				}
			}
			for _, v := range tickermap["bids"].([]any) {
				for i, vv := range v.([]any) {
					switch i {
					case 0:
						ticker.Bid, _ = strconv.ParseFloat(vv.(string), 64)
					case 1:
						ticker.BidSize = vv.(float64)
					}
				}
			}
			ticker.Pair = pairrev
			if ticker.Valid() {
				mdt.Store(pair, ticker)
			}
		} else if isFrozen != "0" {
			// logs.I(pair, " isFrozen polo ")
		} else {
			logs.I(" polo ", pair, isFrozen, o1)
		}
	}
	return
}

// func (polo *Poloniex) OrderBook() (mdt *sync.Map, err error) {
//
//		url := PUBLIC_URL + TICKER_API
//		ret, err := HttpGet(polo.client, url)
//		if err != nil {
//			return nil, err
//		}
//		for pairrev, v := range ret {
//			pair := Sym2duo(pairrev)
//			tickermap, ok := v.(map[string]any)
//			if !ok {
//				return nil, errors.New("not found")
//			}
//			isFrozen, _ := strconv.ParseFloat(tickermap["isFrozen"].(string), 64)
//			if isFrozen == 0 {
//				var ticker BBO
//				ticker.Pair = pairrev
//				ticker.Bid, _ = strconv.ParseFloat(tickermap["highestBid"].(string), 64)
//				ticker.Ask, _ = strconv.ParseFloat(tickermap["lowestAsk"].(string), 64)
//				ticker.BidSize = 1e-12
//				ticker.AskSize = 1e-12
//				ticker.High, _ = strconv.ParseFloat(tickermap["high24hr"].(string), 64)
//				ticker.Low, _ = strconv.ParseFloat(tickermap["low24hr"].(string), 64)
//				ticker.Last, _ = strconv.ParseFloat(tickermap["last"].(string), 64)
//				ticker.Vol, _ = strconv.ParseFloat(tickermap["quoteVolume"].(string), 64)
//				if ticker.Valid() {
//					mdt.Store(pair,ticker)
//				}
//			} else {
//				logs.I(pair, "isFrozen")
//			}
//		}
//		return
//	}
func (poloniex *Poloniex) Post(params url.Values) (resp []byte, err error) {
	sign, err := poloniex.buildPostForm(&params)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"Key":  poloniex.accessKey,
		"Sign": sign}
	resp, err = HttpPostForm2(poloniex.client, TRADE_API, params, headers)
	if err != nil {
		logs.E(err)
		return nil, err
	}
	return resp, err
}
