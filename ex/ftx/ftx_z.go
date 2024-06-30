package ftx

import (
	"github.com/conbanwa/wstrader/q"
	"strings"
	"sync"
	"time"

	"github.com/conbanwa/logs"
)

type Symbol struct {
	Success bool `json:"success"`
	Result  []struct {
		Name                  string  `json:"name"`
		Enabled               bool    `json:"enabled"`
		PostOnly              bool    `json:"postOnly"`
		PriceIncrement        float64 `json:"priceIncrement"`
		SizeIncrement         float64 `json:"sizeIncrement"`
		MinProvideSize        float64 `json:"minProvideSize"`
		Last                  float64 `json:"last"`
		Bid                   float64 `json:"bid"`
		Ask                   float64 `json:"ask"`
		Price                 float64 `json:"price"`
		Type                  string  `json:"type"`
		BaseCurrency          string  `json:"baseCurrency"`
		QuoteCurrency         string  `json:"quoteCurrency"`
		Underlying            string  `json:"underlying"`
		Restricted            bool    `json:"restricted"`
		HighLeverageFeeExempt bool    `json:"highLeverageFeeExempt"`
		Change1H              float64 `json:"change1h"`
		Change24H             float64 `json:"change24h"`
		ChangeBod             float64 `json:"changeBod"`
		QuoteVolume24H        float64 `json:"quoteVolume24h"`
		VolumeUsd24H          float64 `json:"volumeUsd24h"`
		TokenizedEquity       bool    `json:"tokenizedEquity,omitempty"`
	} `json:"result"`
}

func (client *Client) PairArray() (map[string]q.D, map[q.D]q.P, error) {
	s2d := map[string]q.D{}
	d2p := map[q.D]q.P{}
	resp, err := client._get("markets", []byte(""))
	if err != nil {
		return s2d, d2p, err
	}
	var atks Symbol
	err = _processResponse(resp, &atks)
	if err != nil {
		return s2d, d2p, err
	}
	for _, v := range atks.Result {
		if v.Enabled && v.Type == "spot" && !v.Restricted {
			d := q.D{Base: unify(v.BaseCurrency), Quote: unify(v.QuoteCurrency)}
			s2d[v.Name] = d
			var prec q.P
			prec.Base = v.SizeIncrement
			prec.Quote = v.PriceIncrement
			prec.Price = v.PriceIncrement
			prec.MinBase = v.MinProvideSize
			// prec.MinQuote = v.MinValue
			if d.Valid() {
				d2p[d] = prec
			} else {
				logs.E(v)
			}
		}
	}
	return s2d, d2p, err
}
func Sym2duo(pair string) q.D {
	parts := strings.Split(pair, "_")
	if len(parts) != 2 {
		panic("FATAL: SPLIT ERR! " + pair)
	}
	if parts[0] == parts[1] {
		panic("FATAL: SAME CURRENCY ERR! "+ pair)
	}
	return q.D{Base: unify(parts[0]), Quote: unify(parts[1])}
}
func (client *Client) Fee() float64 {
	return 0.0007
}
func (client *Client) PlaceOrders(places [3]q.Order) (orders [3]q.Order, err error) {
	orders = places
	for _, p := range places {
		p.DealAmount = p.Amount
		// go func(p Order) {
		// }(p)
	}
	time.Sleep(time.Second)
	return
}
func (client *Client) Balances() (availables, frozens *sync.Map, err error) {
	availables, frozens = new(sync.Map), new(sync.Map)
	// subaccounts, err := client.GetSubaccountBalances(LILY)
	// zelo.PanicOnErr(err).Send()
	return
}
func (client *Client) Test() bool {
	return true
}
func unify(local string) string {
	global := strings.ToUpper(local)
	return global
}
func (client *Client) GetAttr() (a q.Attr) {
	a.MultiOrder = false
	return a
}

func (client *Client) FeeMap() (map[string]q.TradeFee, error) {
	return nil, nil
}
func (client *Client) WithdrawFee() (sf []q.NetworkWithdraw, err error) {

	return
}
func (client *Client) OneTicker(d q.D) (ticker q.Bbo, err error) {
	return
}
func (client *Client) AllTicker(SymPair map[string]q.D) (mdt *sync.Map, err error) {
	resp, err := client._get("markets", []byte(""))
	if err != nil {
		return
	}
	var atks Symbol
	err = _processResponse(resp, &atks)
	if err != nil {
		return
	}
	for _, v := range atks.Result {
		if v.Enabled && v.Type == "spot" && !v.Restricted {
			d := q.D{Base: unify(v.BaseCurrency), Quote: unify(v.QuoteCurrency)}
			var ticker q.Bbo
			ticker.Pair = v.Name
			ticker.Bid = v.Bid
			ticker.BidSize = v.MinProvideSize
			ticker.Ask = v.Ask
			ticker.AskSize = v.MinProvideSize
			if ticker.Valid() {
				mdt.Store(d, ticker)
			}
		}
	}
	return
}
