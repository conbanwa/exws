package okex

import (
	"math"
	"qa3/wstrader/q"
	"strings"
	"sync"

	"github.com/conbanwa/logs"
)

const PairUrl = baseUrl + "/api/spot/v3/instruments"
const TickerUrl = PairUrl + "/ticker"

func (ok *OKEx) PairArray() (map[string]q.D, map[q.D]q.P, error) {
	// var response []struct {
	// 	InstrumentId  string  `json:"instrument_id"`
	// 	BaseCurrency  string  `json:"base_currency"`
	// 	QuoteCurrency string  `json:"quote_currency"`
	// 	MinSize       float64 `json:"min_size,string"`
	// 	SizeIncrement string  `json:"size_increment"`
	// 	TickSize      string  `json:"tick_size"`
	// }
	var sm = map[string]q.D{}
	var ps = map[q.D]q.P{}
	response, err := ok.OKExSpot.GetCurrenciesPrecision()
	if err != nil {
		return nil, nil, err
	}
	for _, tickermap := range response {
		sm[tickermap.Symbol] = Sym2duo(tickermap.Symbol)
		ps[sm[tickermap.Symbol]] = q.P{Base: math.Pow10(-tickermap.AmountPrecision)}
	}
	return sm, ps, nil
}
func Sym2duo(pair string) q.D {
	parts := strings.Split(pair, "-")
	var res q.D
	if len(parts) == 2 {
		res = q.D{Base: parts[0], Quote: parts[1]}
	} else {
		logs.F("FATAL: DIV ERR!", pair)
	}
	return res
}
func (ok *OKEx) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	orders = places
	for _, p := range places {
		go func(p q.Order) {
		}(p)
	}
	return
}
func (ok *OKEx) Test() bool {
	return true
}
func unify(local string) string {
	global := strings.ToUpper(local)
	return global
}
func (ok *OKEx) GetAttr() (a q.Attr) {
	a.MultiOrder = false
	return a
}

func (ok *OKEx) OneTicker(d q.D) (ticker q.Bbo, err error) {
	return
}
func (ok *OKEx) AllTicker(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	var response []spotTickerResponse
	err = ok.DoRequest("GET", "/api/spot/v3/instruments/ticker", "", &response)
	if err != nil {
		return
	}
	for _, tickermap := range response {
		pair := Sym2duo(tickermap.InstrumentId)
		var ticker q.Bbo
		ticker.Pair = tickermap.InstrumentId
		ticker.Bid = tickermap.BestBid
		ticker.Ask = tickermap.BestAsk
		ticker.BidSize = tickermap.BidSize
		ticker.AskSize = tickermap.AskSize
		// ticker.High = tickermap.High24h
		// ticker.Low = tickermap.Low24h
		// ticker.Last = tickermap.Last
		// ticker.Vol = tickermap.BaseVolume24h
		// ticker.Date = 0
		if ticker.Valid() {
			mdt.Store(pair, ticker)
		}
	}
	return
}
